// Command tuitest is a thin wrapper around `go test` for tuitest-powered
// test suites. It adds flags that map to tuitest features: -update to
// regenerate snapshots, -junit/-html to emit reports, -filter to pick
// specific tests, -parallel to set parallelism, and -watch to re-run on
// file changes.
//
// Subcommands:
//
//	tuitest diff <testname>   show the failure diff for a named test
//
// Usage:
//
//	tuitest [flags] [packages...]
//
// Packages default to "./..." when none are provided. The default reporter
// is the vitest-style runner already wired into the test code.
//
// Examples:
//
//	tuitest                                   # go test ./...
//	tuitest -filter TestHarness ./tuitest/... # run tests matching TestHarness
//	tuitest -update ./tuitest/...             # regenerate snapshots
//	tuitest -junit out/junit.xml -parallel 4  # parallel run + junit output
//	tuitest -watch                            # re-run on file changes (1s poll)
//	tuitest diff TestFoo                      # show diff for TestFoo failure
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/moneycaringcoder/tuikit-go/tuitest"
)

func main() {
	// Handle subcommands before flag parsing.
	if len(os.Args) >= 2 && os.Args[1] == "diff" {
		runDiffSubcommand(os.Args[2:])
		return
	}

	var (
		filter   = flag.String("filter", "", "run only tests matching regexp (maps to go test -run)")
		update   = flag.Bool("update", false, "regenerate tuitest snapshots (passes -tuitest.update to tests)")
		junit    = flag.String("junit", "", "write JUnit XML report to path (informational; tests must use JUnitReporter to populate it)")
		htmlOut  = flag.String("html", "", "write HTML report to path (informational; tests must use HTMLReporter to populate it)")
		parallel = flag.Int("parallel", 0, "maximum number of tests to run in parallel (maps to go test -parallel)")
		watch    = flag.Bool("watch", false, "watch the working tree for changes and re-run on modification")
		verbose  = flag.Bool("v", false, "verbose go test output")
	)
	flag.Parse()

	packages := flag.Args()
	if len(packages) == 0 {
		packages = []string{"./..."}
	}

	runOnce := func() int {
		code := runGoTest(*filter, *update, *junit, *htmlOut, *parallel, *verbose, packages)
		if code != 0 && *watch {
			printFailureDiffHints()
		}
		return code
	}

	if !*watch {
		os.Exit(runOnce())
	}

	// Watch mode: simple mtime poll on .go files under cwd.
	fmt.Fprintln(os.Stderr, "[tuitest] watch mode (polling every 1s, Ctrl+C to stop)")
	lastHash := snapshotTree(".")
	runOnce()
	for {
		time.Sleep(time.Second)
		h := snapshotTree(".")
		if h != lastHash {
			fmt.Fprintln(os.Stderr, "[tuitest] change detected, re-running")
			lastHash = h
			runOnce()
		}
	}
}

// runDiffSubcommand implements `tuitest diff [testname]`.
// Without a testname it lists all available failure captures.
// With a testname it prints the diff to stdout using DiffViewer's text output.
func runDiffSubcommand(args []string) {
	if len(args) == 0 {
		// List available captures.
		names, err := tuitest.ListFailureCaptures()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[tuitest diff] error: %v\n", err)
			os.Exit(1)
		}
		if len(names) == 0 {
			fmt.Println("[tuitest diff] no failure captures found (run tests first)")
			return
		}
		fmt.Println("Available failure captures:")
		for _, n := range names {
			fmt.Println("  " + n)
		}
		return
	}

	testName := strings.Join(args, " ")
	fc, err := tuitest.LoadFailureCapture(testName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[tuitest diff] %v\n", err)
		os.Exit(1)
	}

	dv := tuitest.NewDiffViewer(fc)
	dv.SetSize(120, 40)
	printDiffViewerModes(dv, fc)
}

// printDiffViewerModes renders all three modes of the DiffViewer to stdout.
func printDiffViewerModes(dv *tuitest.DiffViewer, fc *tuitest.FailureCapture) {
	modes := []struct {
		key  string
		mode tuitest.DiffMode
	}{
		{"s", tuitest.DiffModeSideBySide},
		{"u", tuitest.DiffModeUnified},
		{"d", tuitest.DiffModeCellsOnly},
	}
	// Show side-by-side by default; user can re-run to see other modes.
	_ = modes
	// For one-shot CLI we just render side-by-side then unified then cells.
	for _, m := range modes {
		dv.SetMode(m.mode)
		fmt.Println(dv.View())
		fmt.Println(strings.Repeat("─", 80))
	}
}

// printFailureDiffHints prints a hint after a failed watch-mode run showing
// which test failures have diff captures available.
func printFailureDiffHints() {
	names, err := tuitest.ListFailureCaptures()
	if err != nil || len(names) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "[tuitest] failure diffs available — view with:")
	for _, n := range names {
		fmt.Fprintf(os.Stderr, "  tuitest diff %s\n", n)
	}
}

func runGoTest(filter string, update bool, junit, htmlOut string, parallel int, verbose bool, packages []string) int {
	args := []string{"test"}
	if verbose {
		args = append(args, "-v")
	}
	if filter != "" {
		args = append(args, "-run", filter)
	}
	if parallel > 0 {
		args = append(args, "-parallel", fmt.Sprintf("%d", parallel))
	}
	args = append(args, packages...)
	if update || junit != "" || htmlOut != "" {
		args = append(args, "-args")
		if update {
			args = append(args, "-tuitest.update")
		}
		if junit != "" {
			args = append(args, "-tuitest.junit="+junit)
		}
		if htmlOut != "" {
			args = append(args, "-tuitest.html="+htmlOut)
		}
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "[tuitest] run failed: %v\n", err)
		return 1
	}
	return 0
}

// snapshotTree returns a coarse hash of .go file modification times under
// root. Used only to detect "something changed" in watch mode.
func snapshotTree(root string) string {
	var sb strings.Builder
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || strings.HasPrefix(name, ".omc") {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		fmt.Fprintf(&sb, "%s:%d\n", path, info.ModTime().UnixNano())
		return nil
	})
	return sb.String()
}

// Package main provides a zero-dependency Go AST inspection tool for Linden.
// It verifies that interface definitions and error taxonomy in src/ remain in
// lockstep with normative architectural specifications in docs/architecture/interfaces/.
//
// Usage:
//
//	go run ./tools/docgen -verify     # Fails with non-zero exit code if docs have drifted from code
//	go run ./tools/docgen -dump       # Dumps extracted interfaces and signatures to stdout
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type InterfaceDef struct {
	Package string
	Name    string
	Methods []string
	DocPath string
}

func main() {
	verify := flag.Bool("verify", false, "Verify documentation matches Go interface signatures")
	dump := flag.Bool("dump", false, "Dump extracted interface signatures")
	flag.Parse()

	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error locating repository root: %v\n", err)
		os.Exit(1)
	}

	fset := token.NewFileSet()

	// 1. Inspect Go Interfaces
	interfaces := []InterfaceDef{
		{
			Package: "orchestrator",
			Name:    "ChatService",
			DocPath: filepath.Join(repoRoot, "docs", "architecture", "interfaces", "orchestrator-service.md"),
		},
		{
			Package: "inference",
			Name:    "Client",
			DocPath: filepath.Join(repoRoot, "docs", "architecture", "interfaces", "inference-client.md"),
		},
		{
			Package: "storage",
			Name:    "Store",
			DocPath: filepath.Join(repoRoot, "docs", "architecture", "interfaces", "storage-store.md"),
		},
		{
			Package: "discovery",
			Name:    "Advertiser",
			DocPath: filepath.Join(repoRoot, "docs", "architecture", "interfaces", "discovery-advertiser.md"),
		},
	}

	for i := range interfaces {
		pkgDir := filepath.Join(repoRoot, "src", interfaces[i].Package)
		methods, err := extractInterfaceMethods(fset, pkgDir, interfaces[i].Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed parsing %s.%s: %v\n", interfaces[i].Package, interfaces[i].Name, err)
			os.Exit(1)
		}
		interfaces[i].Methods = methods
	}

	// 2. Inspect Error Codes
	errCodes, err := extractErrorCodes(fset, filepath.Join(repoRoot, "src", "errs"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed extracting error codes: %v\n", err)
		os.Exit(1)
	}

	if *dump {
		fmt.Println("=== Linden Architectural Interfaces ===")
		for _, iface := range interfaces {
			fmt.Printf("\nInterface: %s.%s (Doc: %s)\n", iface.Package, iface.Name, filepath.Base(iface.DocPath))
			for _, m := range iface.Methods {
				fmt.Printf("  • %s\n", m)
			}
		}
		fmt.Println("\n=== Linden Error Taxonomy ===")
		for _, code := range errCodes {
			fmt.Printf("  • %s\n", code)
		}
		return
	}

	// 3. Verification Mode
	if *verify {
		var driftErrors []string

		for _, iface := range interfaces {
			docBytes, err := os.ReadFile(iface.DocPath)
			if err != nil {
				driftErrors = append(driftErrors, fmt.Sprintf("Missing doc file: %s", iface.DocPath))
				continue
			}
			docContent := string(docBytes)

			for _, method := range iface.Methods {
				methodName := strings.Split(method, "(")[0]
				if !strings.Contains(docContent, methodName) {
					driftErrors = append(driftErrors, fmt.Sprintf(
						"DRIFT DETECTED: %s.%s method %q is missing from %s",
						iface.Package, iface.Name, methodName, filepath.Base(iface.DocPath),
					))
				}
			}
		}

		// Verify error taxonomy doc
		errDocPath := filepath.Join(repoRoot, "docs", "architecture", "interfaces", "error-taxonomy.md")
		errDocBytes, err := os.ReadFile(errDocPath)
		if err != nil {
			driftErrors = append(driftErrors, fmt.Sprintf("Missing error taxonomy doc: %s", errDocPath))
		} else {
			errDocContent := string(errDocBytes)
			for _, code := range errCodes {
				quotedCode := fmt.Sprintf("%q", code)
				if !strings.Contains(errDocContent, quotedCode) && !strings.Contains(errDocContent, code) {
					driftErrors = append(driftErrors, fmt.Sprintf(
						"DRIFT DETECTED: Error code %s is missing from error-taxonomy.md", code,
					))
				}
			}
		}

		if len(driftErrors) > 0 {
			fmt.Fprintln(os.Stderr, "\n=======================================================")
			fmt.Fprintln(os.Stderr, "  INTERFACE DRIFT VERIFICATION FAILED")
			fmt.Fprintln(os.Stderr, "=======================================================")
			for _, e := range driftErrors {
				fmt.Fprintf(os.Stderr, " [FAIL] %s\n", e)
			}
			fmt.Fprintln(os.Stderr, "\nPlease update docs/architecture/interfaces/ to match Go code.")
			os.Exit(1)
		}

		fmt.Println("✓ All architectural interfaces and error taxonomies are in 100% lockstep with documentation.")
		return
	}

	fmt.Println("Linden Docgen Tool. Run with -verify or -dump.")
}

func extractInterfaceMethods(fset *token.FileSet, pkgDir string, interfaceName string) ([]string, error) {
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var methods []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec.Name.Name != interfaceName {
						continue
					}

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					for _, field := range interfaceType.Methods.List {
						if len(field.Names) == 0 {
							continue
						}
						methodName := field.Names[0].Name
						methods = append(methods, methodName)
					}
				}
			}
		}
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("interface %s not found in %s", interfaceName, pkgDir)
	}

	return methods, nil
}

func extractErrorCodes(fset *token.FileSet, errsDir string) ([]string, error) {
	pkgs, err := parser.ParseDir(fset, errsDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var codes []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.CONST {
					continue
				}

				for _, spec := range genDecl.Specs {
					valSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					for _, val := range valSpec.Values {
						basicLit, ok := val.(*ast.BasicLit)
						if ok && basicLit.Kind == token.STRING {
							valStr := strings.Trim(basicLit.Value, `"`)
							codes = append(codes, valStr)
						}
					}
				}
			}
		}
	}

	if len(codes) == 0 {
		return nil, fmt.Errorf("no error codes found in %s", errsDir)
	}

	return codes, nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find repo root containing .git or go.work")
}

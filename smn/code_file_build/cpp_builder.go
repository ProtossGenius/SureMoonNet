package code_file_build

//CppImp import cpp pakcage.
func CppImp(pkg string) string {
	if len(pkg) == 0 {
		return ""
	}

	if pkg[0] == '"' || pkg[0] == '<' {
		return "#include " + pkg
	}

	return "#include <" + pkg + ">"
}

func CppPkg(cf *CodeFile, pkg string) {

}

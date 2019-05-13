package pgt_achieve

type Hello interface {
	helloa(aaa string) string
	hellob(aaa, bbb string) string
	helloc(ls ...int)
}

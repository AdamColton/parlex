package reducer

import (
	"testing"

	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
)

func TestCanParseEmpty(t *testing.T) {
	rdcr, err := Parse("")
	assert.NoError(t, err)
	assert.NotNil(t, rdcr)
}

func TestPromoteSingleChild(t *testing.T) {
	str := `
		Foo PromoteSingleChild()
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		A {
			A1 {
				A11: "A11"
				A12: "A12"
			}
		}
		B {
			B1 {
				B11: "B11"
				B12: "B12"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)

	pn1, err = tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got = pn1.String()
	expected = `Root {
	A {
		A1 {
			A11: "A11"
			A12: "A12"
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestPromoteGrandChild(t *testing.T) {
	str := `
		Foo PromoteGrandChildren()
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root{
			Foo{
				A{
					A1{
						A11:"A11"
						A12:"A12"
					}
				}
				B{
					B1{
						B11:"B11"
						B12:"B12"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		A1 {
			A11: "A11"
			A12: "A12"
		}
		B1 {
			B11: "B11"
			B12: "B12"
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestPromoteChildrenOf(t *testing.T) {
	str := `
		Foo PromoteChildrenOf(1)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root{
			Foo{
				A{
					A1{
						A11:"A11"
						A12:"A12"
					}
				}
				B{
					B1{
						B11:"B11"
						B12:"B12"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		A {
			A1 {
				A11: "A11"
				A12: "A12"
			}
		}
		B1 {
			B11: "B11"
			B12: "B12"
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestRemoveAll(t *testing.T) {
	str := `
		Foo RemoveAll("B")
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		A {
			A1 {
				A11: "A11"
				A12: "A12"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestRemoveChildren(t *testing.T) {
	str := `
		Foo RemoveChildren(0,2)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		B {
			B1 {
				B11: "B11"
				B12: "B12"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestPromoteChild(t *testing.T) {
	str := `
		Foo PromoteChild(1)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	B {
		A {
			A1 {
				A11: "A11"
				A12: "A12"
			}
		}
		B1 {
			B11: "B11"
			B12: "B12"
		}
		B {
			B2 {
				B21: "B21"
				B22: "B22"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestRemoveChild(t *testing.T) {
	str := `
		Foo RemoveChild(1)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo {
		A {
			A1 {
				A11: "A11"
				A12: "A12"
			}
		}
		B {
			B2 {
				B21: "B21"
				B22: "B22"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestReplaceWithChild(t *testing.T) {
	str := `
		Foo ReplaceWithChild(1)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	B {
		B1 {
			B11: "B11"
			B12: "B12"
		}
	}
}
`
	assert.Equal(t, expected, got)
}

func TestPromoteChildValue(t *testing.T) {
	str := `
		Foo PromoteChildValue(1)
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A: "A"
				B1: "B1"
				B2: "B2" 
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	Foo: "B1" {
		A: "A"
		B2: "B2"
	}
}
`
	assert.Equal(t, expected, got)
}

func TestChildIs(t *testing.T) {
	str := `
		Foo If( ChildIs(0, "A"), PromoteChild(0), PromoteChild(1))
	`
	rdcr := Must(str)

	pn1, err := tree.New(`
		Root {
			Foo {
				A {
					A1 {
						A11: "A11"
						A12: "A12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)

	got := pn1.String()
	expected := `Root {
	A {
		A1 {
			A11: "A11"
			A12: "A12"
		}
		B {
			B1 {
				B11: "B11"
				B12: "B12"
			}
		}
		B {
			B2 {
				B21: "B21"
				B22: "B22"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)

	pn1, err = tree.New(`
		Root {
			Foo {
				C {
					C1 {
						C11: "C11"
						C12: "C12"
					}
				}
				B {
					B1 {
						B11: "B11"
						B12: "B12"
					}
				}
				B {
					B2 {
						B21: "B21"
						B22: "B22"
					}
				}
			}
		}
	`)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)
	assert.Equal(t, expected, got)

	got = pn1.String()
	expected = `Root {
	B {
		C {
			C1 {
				C11: "C11"
				C12: "C12"
			}
		}
		B1 {
			B11: "B11"
			B12: "B12"
		}
		B {
			B2 {
				B21: "B21"
				B22: "B22"
			}
		}
	}
}
`
	assert.Equal(t, expected, got)
}

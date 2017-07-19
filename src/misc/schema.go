package misc

import "github.com/skycoin/cxo/skyobject"

type SchemaRefs struct {
	Board        skyobject.SchemaReference
	Thread       skyobject.SchemaReference
	Post         skyobject.SchemaReference
	ThreadPage   skyobject.SchemaReference
	BoardPage    skyobject.SchemaReference
	ExternalRoot skyobject.SchemaReference
}

func GetSchemaRefsFromRoot(r *skyobject.Root) (*SchemaRefs, error) {
	sr := new(SchemaRefs)
	reg, e := r.Registry()
	if e != nil {
		return nil, e
	}
	sr.Board, e = reg.SchemaReferenceByName("Board")
	if e != nil {
		return nil, e
	}
	sr.Thread, e = reg.SchemaReferenceByName("Thread")
	if e != nil {
		return nil, e
	}
	sr.Post, e = reg.SchemaReferenceByName("Post")
	if e != nil {
		return nil, e
	}
	sr.ThreadPage, e = reg.SchemaReferenceByName("ThreadPage")
	if e != nil {
		return nil, e
	}
	sr.BoardPage, e = reg.SchemaReferenceByName("BoardPage")
	if e != nil {
		return nil, e
	}
	sr.ExternalRoot, e = reg.SchemaReferenceByName("ExternalRoot")
	if e != nil {
		return nil, e
	}
	return sr, nil
}

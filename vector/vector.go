package vector

import (
	"errors"
	"reflect"
	"study/simple-stl"
)

var (
	rtw = reflect.TypeOf((*simple_stl.Writer)(nil)).Elem()
	rtr = reflect.TypeOf((*simple_stl.Reader)(nil)).Elem()
)

type Vector interface {
	At(int) (interface{})
	Edit(int, interface{})
	Append(writer interface{})
}

type vect struct {
	simple_stl.Sdata
}

func NewSdata(model interface{}) Vector {
	return &vect{
		Sdata: simple_stl.NewSdata(model),
	}
}

func (t *vect) At(i int) interface{} {
	if len(t.Sdata.Data) < i+1 {
		t.Err = errors.New("index out of range")
		return nil
	}
	cw,err:= permissionCheck(t.Model, t.Sdata.Data[i], rtr)
	if err != nil{
		t.Err = err
		return nil
	}
	t.Sdata.Rwm.Lock()
	defer t.Sdata.Rwm.Unlock()
	return cw.(simple_stl.Reader).Read()

}

func (t *vect) Edit(i int, w interface{}) {
	if len(t.Sdata.Data) < i+1 {
		panic("index out of range")
		return
	}

	cw,err:= permissionCheck(t.Model, w, rtw)
	if err != nil{
		t.Err = err
		return
	}
	t.Sdata.Rwm.Lock()
	defer t.Sdata.Rwm.Unlock()
	t.Sdata.Data[i].(simple_stl.Writer).Write(cw.(simple_stl.Writer))

}

func (t *vect) Append(w interface{}) {
	cw,err:=permissionCheck(t.Model, w, rtw)
	if err != nil{
		t.Err = err
		return
	}
	t.Sdata.Rwm.Lock()
	defer t.Sdata.Rwm.Unlock()
	t.Sdata.Data = append(t.Sdata.Data, cw)
}

func permissionCheck(p reflect.Type, w interface{}, chk reflect.Type) (interface{}, error){
	if !p.Implements(chk){
		return nil, errors.New("permission denied")
	}
	if reflect.TypeOf(w).Implements(chk){
		return w, nil
	}
	if pw:=transform(w); reflect.TypeOf(pw).Implements(chk){
		return pw, nil
	}
	return nil, errors.New(reflect.TypeOf(w).String() + " is not " + chk.String())
}

func transform( w interface{}) interface{} {
	pv := reflect.New(reflect.TypeOf(w))
	pv.Elem().Set(reflect.ValueOf(w))
	return pv.Interface()
}

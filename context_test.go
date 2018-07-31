package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"github.com/valyala/fasthttp"
	"errors"
	"fmt"
)

type ErrornousJson struct {
}

func (*ErrornousJson) MarshalJSON() ([]byte, error) {
	return nil, errors.New("forced error")
}

var _ = g.Describe("Http", func() {
	g.Describe("Context", func() {
		g.It("should accept the Content-type as JSON", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			ctx.Request.Header.SetContentType("application/json")
			ctx.Request.AppendBodyString(`{"foo":"bar"}`)
			Expect(ctx.IsJson()).To(BeTrue())
		})

		g.It("should accept the Content-type with charset as JSON", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			ctx.Request.Header.SetContentType("application/json; charset=utf-8")
			ctx.Request.AppendBodyString(`{"foo":"bar"}`)
			Expect(ctx.IsJson()).To(BeTrue())
		})

		g.It("should not accept the Content-type as JSON", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			Expect(ctx.IsJson()).To(BeFalse())
		})

		g.It("should return the JSON structure in the body of the requisition", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			ctx.Request.Header.SetContentType("application/json; charset=utf-8")
			ctx.Request.AppendBodyString(`{"foo":"bar"}`)
			var data struct {
				Foo string `json:"foo"`
			}
			Expect(ctx.BodyJson(&data)).To(BeNil())
			Expect(data.Foo).To(Equal("bar"))
		})

		g.It("should return some user values defined", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			ctx.Ctx.SetUserValue("key1", "value1")
			ctx.Ctx.SetUserValue("key2", 2)
			Expect(ctx.UserValue("key1")).To(Equal("value1"))
			Expect(ctx.UserValue("key2")).To(Equal(2))
		})

		g.It("should serialize and send a JSON with the rightful 'Content-type' header", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			Expect(ctx.SendJson(map[string]interface{}{
				"foo": "bar",
			})).To(BeNil())
			Expect(string(ctx.Response.Header.ContentType())).To(Equal("application/json"))
			tmp := bytes.NewBufferString("")
			ctx.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`{"foo":"bar"}`))
		})

		g.It("should fail serializing a JSON struct", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			Expect(ctx.SendJson(&ErrornousJson{})).NotTo(BeNil())
		})

		g.It("should write a buff to the context body", func() {
			ctx := NewContext(&fasthttp.RequestCtx{})
			n, err := fmt.Fprint(ctx, "this is a test")
			Expect(err).To(BeNil())
			Expect(n).To(Equal(14))
			tmp := bytes.NewBufferString("")
			ctx.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})
	})
})
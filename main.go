package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

func main() {
	publishableKey := "pk_test_5BV29MJdeaXwC6UpVHeY8uQM00VTvLqATm"
	stripe.Key = "sk_test_zhVLqm2IiBSzHCmYZbJwlwB400fY8QkLs2"

	p := pat.New()

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, map[string]string{"Key": publishableKey})
	})

	p.Post("/charge", func(res http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		email := req.Form.Get("stripeEmail")

		customerParams := &stripe.CustomerParams{Email: &email}
		customerParams.SetSource(req.Form.Get("stripeToken"))

		newCustomer, err := customer.New(customerParams)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		var amount int64
		amount = 500
		value := "usd"
		id := newCustomer.ID
		disc := "Sample Charge"
		chargeParams := &stripe.ChargeParams{
			Amount:      &amount,
			Currency:    &value,
			Description: &disc,
			Customer:    &id,
		}

		if _, err := charge.New(chargeParams); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(res, "Charge completed successfully!")

	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":80", p))
}

var indexTemplate = `
<html>
<head>
    <title>Checkout Example</title>
</head>
<body>
    <form action="/charge" method="post" class="payment">
        <article>
          <label class="amount">
            <span>Amount: $5.00</span>
          </label>
        </article>
    
        <script src="https://checkout.stripe.com/checkout.js" class="stripe-button"
                data-key="{{ .Key }}"
                data-description="A month's subscription"
                data-amount="500"
                data-locale="auto"></script>
    </form>
</body>
</html>
`

package graphql

func schema() {
	/*
		userType := graphql.NewObject(graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id":   &graphql.Field{Type: graphql.Int},
				"name": &graphql.Field{Type: graphql.String},
			},
		})

		rootQuery := graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: userType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						user := map[string]interface{}{
							"id":   1,
							"name": "Jera",
						}
						return user, nil
					},
				},
			},
		})

		//schema, err := graphql.NewSchema(graphql.SchemaConfig{
		//	Query: rootQuery,
		//})
		//if err != nil {
		//	panic(err)
			router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
				result := graphql.Do(graphql.Params{
					Schema:        schema,
					RequestString: r.URL.Query().Get("query"),
				})
				json.NewEncoder(w).Encode(result)
			}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	*/
}

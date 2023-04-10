docker-compose up -d
to run application



    _______________________________________________________________________________________________
	app.Post("/Auth/Token", controllers.CreateToken)//login
    _______________________________________________________________________________________________

    app.Post("/User/Register", controllers.Register)/registration

	app.Get("/Users", middleware.JWTProtected(), controllers.GetUsers)//admin get all users

	app.Delete("/User/Delete/:userId", middleware.JWTProtected(), controllers.DeleteUser)//delete
    _______________________________________________________________________________________________

    app.Get("/Visit/GainQr", middleware.JWTProtected(), controllers.CreateQr)/create qr(terminal)

	app.Get("/Visits", middleware.JWTProtected(), controllers.GetVisits) get all visits(Mobilka)

	app.Delete("/Visits/DeleteLegasy", middleware.JWTProtected(), controllers.DeleteLegacyVisits)//delete visits admin mobile

	app.Get("/Visits/:code", middleware.JWTProtected(), controllers.CreateVisit)//create visit scan qr 
    _______________________________________________________________________________________________
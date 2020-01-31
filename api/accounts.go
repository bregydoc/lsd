package api

func (api *API) registerNewAdminAccount(username, password string) {
	if api.accounts == nil {
		api.accounts = map[string]string{}
	}
	api.accounts[username] = password
}

func (api *API) RegisterAccount(username, password string) {
	api.registerNewAdminAccount(username, password)
}

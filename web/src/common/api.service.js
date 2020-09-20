/** @format */

import axios from "axios";

const ApiService = {
	getURL(endpoint) {
		return process.env.VUE_APP_TOWER_URL.replace(/\/$/, "") + endpoint;
	},

	getHeaders() {
		let apiKey = "";

		if (localStorage.getItem("user_api_token") != null) {
			apiKey = localStorage.getItem("user_api_token");
		}

		return {
			crossdomain: true,
			headers: {
				"x-api-key": apiKey,
				"Content-Type": "application/json",
			},
		};
	},

	get(endpoint) {
		return axios.get(this.getURL(endpoint), this.getHeaders());
	},

	delete(endpoint) {
		return axios.delete(this.getURL(endpoint), this.getHeaders());
	},

	post(endpoint, data = {}) {
		return axios.post(this.getURL(endpoint), data, this.getHeaders());
	},

	put(endpoint, data = {}) {
		return axios.put(this.getURL(endpoint), data, this.getHeaders());
	},
};

export default ApiService;

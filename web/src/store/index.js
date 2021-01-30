/** @format */

import Vue from "vue";
import Vuex from "vuex";
import tower from "./tower.module";
import auth from "./auth.module";
import settings from "./settings.module";
import host from "./host.module";

Vue.use(Vuex);

export default new Vuex.Store({
	modules: {
		tower,
		auth,
		settings,
		host,
	},
});

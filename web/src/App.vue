<!-- @format -->

<template>
	<div id="app">
		<div id="nav">
			<router-link to="/"
				><b-icon pack="fas" icon="home" size="is-small"> </b-icon>
				Home</router-link
			>
			<template v-if="logged">
				|
				<router-link to="/hosts">
					<b-icon pack="fas" icon="server" size="is-small"> </b-icon>
					Hosts</router-link
				>
				|
				<router-link to="/jobs">
					<b-icon pack="fas" icon="bolt" size="is-small"> </b-icon>
					Jobs</router-link
				>
				|
				<router-link to="/settings">
					<b-icon pack="fas" icon="cog" size="is-small"> </b-icon>
					Settings</router-link
				>
				|
				<a href="#" @click="logout">
					<b-icon pack="fas" icon="sign-out-alt" size="is-small">
					</b-icon>
					Logout</a
				>
			</template>
			<template v-else>
				|
				<router-link to="/login">
					<b-icon pack="fas" icon="sign-in-alt" size="is-small">
					</b-icon>
					Login</router-link
				>
			</template>
		</div>
		<router-view @refresh-state="refreshState" />
	</div>
</template>

<style>
#app {
	text-align: center;
	color: #2c3e50;
}

#nav {
	padding: 30px;
}

#nav a {
	font-weight: bold;
	color: #2c3e50;
}

#nav a.router-link-exact-active {
	color: #42b983;
}
</style>

<script>
export default {
	data() {
		return {
			logged: localStorage.getItem("user_api_token") != null,
		};
	},
	methods: {
		logout() {
			console.log("Logout");
			this.logged = false;
			localStorage.removeItem("user_api_token");
			this.$router.push("/login");
		},
		refreshState() {
			this.logged = localStorage.getItem("user_api_token") != null;
		},
	},
	mounted() {
		this.logged = localStorage.getItem("user_api_token") != null;
	},
};
</script>

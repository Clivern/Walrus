<!-- @format -->

<template>
	<div class="columns is-mobile is-centered">
		<div class="column is-4">
			<b-notification
				v-show="notification.message"
				:type="notification.type"
				has-icon
				aria-close-label="Close"
			>
				{{ notification.message }}
			</b-notification>

			<div class="field">
				<p class="control has-icons-left has-icons-right">
					<input class="input" type="email" placeholder="Email" />
					<span class="icon is-small is-left">
						<i class="fas fa-envelope"></i>
					</span>
				</p>
			</div>
			<div class="field">
				<p class="control has-icons-left">
					<input
						class="input"
						type="password"
						placeholder="Password"
					/>
					<span class="icon is-small is-left">
						<i class="fas fa-lock"></i>
					</span>
				</p>
			</div>
			<div class="field">
				<p class="control">
					<template v-if="!setup">
						<b-button type="is-success" @click="loginEvent"
							>Login</b-button
						>
					</template>
					<template v-else>
						<b-button type="is-success" @click="setupEvent"
							>Create admin account</b-button
						>
					</template>
				</p>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: "login",

	data() {
		return {
			setup: false,
			notification: {
				type: null,
				message: null,
			},
		};
	},

	methods: {
		loginEvent() {
			console.log(this.$store.getters["auth/getAuthResult"]);

			this.$store.dispatch("auth/authAction", ["email", "password"]).then(
				() => {
					console.log("BAM!");
					console.log(this.getAuthResult);
					localStorage.setItem("user_api_token", "xxxxx");
					this.$router.push("/");
				},
				(err) => {
					this.$buefy.toast.open({
						message: err,
						type: "is-danger",
					});
				}
			);
		},
		setupEvent() {
			console.log("Setup Action");

			this.$store
				.dispatch("auth/setupAction", ["email", "password"])
				.then(
					() => {
						console.log("BAM!");
					},
					(err) => {
						this.$buefy.toast.open({
							message: err,
							type: "is-danger",
						});
					}
				);
		},
	},

	mounted() {
		this.$store.dispatch("auth/fetchInfo").then(
			() => {
				this.setup = this.$store.getters[
					"auth/getTowerInfo"
				].setupStatus;

				if (this.setup) {
					this.notification.type = "is-info";
					this.notification.message =
						"Walrus admin user not created yet, Please submit your desired email and password below!";
				}
			},
			(err) => {
				this.$buefy.toast.open({
					message: err,
					type: "is-danger",
				});
			}
		);
	},
};
</script>

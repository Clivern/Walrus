<!-- @format -->

<template>
	<div class="columns is-desktop is-centered">
		<div class="column is-4">
			<template v-if="setup">
				<div class="field">
					<p class="control has-icons-left has-icons-right">
						<input
							class="input"
							type="text"
							placeholder="Name"
							v-model="form.user_name"
							required
						/>
						<span class="icon is-small is-left">
							<i class="fas fa-user"></i>
						</span>
					</p>
				</div>
			</template>
			<div class="field">
				<p class="control has-icons-left has-icons-right">
					<input
						class="input"
						type="email"
						placeholder="Email"
						v-model="form.user_email"
						required
					/>
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
						v-model="form.user_password"
						required
					/>
					<span class="icon is-small is-left">
						<i class="fas fa-lock"></i>
					</span>
				</p>
			</div>
			<div class="field">
				<p class="control">
					<template v-if="!setup">
						<b-button
							type="submit is-danger is-light"
							v-bind:disabled="form.button_disabled"
							@click="loginEvent"
							>Login</b-button
						>
					</template>
					<template v-else>
						<b-button
							type="is-danger is-light"
							v-bind:disabled="form.button_disabled"
							@click="setupEvent"
							>Install Walrus</b-button
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

			form: {
				user_password: "",
				user_email: "",
				user_name: "",
				button_disabled: false,
			},

			// Loader
			loader: {
				isFullPage: true,
				ref: null,
			},
		};
	},

	methods: {
		loading() {
			this.loader.ref = this.$buefy.loading.open({
				container: this.loader.isFullPage ? null : this.$refs.element.$el,
			});
		},

		loginEvent() {
			this.form.button_disabled = true;

			this.$store
				.dispatch("auth/authAction", {
					email: this.form.user_email,
					password: this.form.user_password,
				})
				.then(
					(response) => {
						this.$buefy.toast.open({
							message: "User logged in successfully",
							type: "is-success",
						});

						localStorage.setItem("user_api_key", response.data.apiKey);

						localStorage.setItem("user_email", response.data.email);
						localStorage.setItem("user_id", response.data.id);
						localStorage.setItem("user_name", response.data.name);
						this.$router.push("/");
					},
					(err) => {
						if (err.response.data.errorMessage) {
							this.$buefy.toast.open({
								message: err.response.data.errorMessage,
								type: "is-danger",
							});
						} else {
							this.$buefy.toast.open({
								message: "Error status code: " + err.response.status,
								type: "is-danger",
							});
						}
						this.form.button_disabled = false;
					}
				);
		},
		setupEvent() {
			this.form.button_disabled = true;

			this.$store
				.dispatch("auth/setupAction", {
					name: this.form.user_name,
					email: this.form.user_email,
					password: this.form.user_password,
				})
				.then(
					() => {
						this.$buefy.toast.open({
							message: "Walrus installed successfully",
							type: "is-success",
						});

						this.$router.push("/");
					},
					(err) => {
						if (err.response.data.errorMessage) {
							this.$buefy.toast.open({
								message: err.response.data.errorMessage,
								type: "is-danger is-light",
							});
						} else {
							this.$buefy.toast.open({
								message: "Error status code: " + err.response.status,
								type: "is-danger is-light",
							});
						}
						this.form.button_disabled = false;
					}
				);
		},
	},

	mounted() {
		this.loading();

		this.$store.dispatch("auth/fetchInfo").then(
			() => {
				let info = this.$store.getters["auth/getTowerInfo"];

				if (!info.setupStatus) {
					this.setup = true;

					this.$buefy.toast.open({
						message: "Walrus not installed yet, Please setup the application!",
						type: "is-info is-light",
					});
				}

				this.loader.ref.close();
			},
			(err) => {
				this.$buefy.toast.open({
					message: err.response.data.errorMessage,
					type: "is-danger is-light",
				});

				this.loader.ref.close();
			}
		);
	},
};
</script>

<!-- @format -->

<template>
	<div class="columns is-desktop is-centered">
		<div class="column is-4">
			<br /><br />
			<section>
				<b-field label="S3 Key">
					<b-input type="text" v-model="form.backup_s3_key"> </b-input>
				</b-field>
				<b-field label="S3 Secret">
					<b-input type="text" v-model="form.backup_s3_secret"> </b-input>
				</b-field>
				<b-field label="S3 Endpoint">
					<b-input type="text" v-model="form.backup_s3_endpoint"> </b-input>
				</b-field>
				<b-field label="S3 Region">
					<b-input type="text" v-model="form.backup_s3_region"> </b-input>
				</b-field>
				<b-field label="S3 Bucket">
					<b-input type="text" v-model="form.backup_s3_bucket"> </b-input>
				</b-field>
				<br />
				<div class="field">
					<p class="control">
						<b-button
							type="is-danger is-light"
							v-bind:disabled="form.button_disabled"
							@click="updateEvent"
							>Update</b-button
						>
					</p>
				</div>
			</section>
		</div>
	</div>
</template>

<script>
export default {
	name: "settings",

	data() {
		return {
			form: {
				button_disabled: false,

				backup_s3_key: "",
				backup_s3_secret: "",
				backup_s3_endpoint: "",
				backup_s3_region: "",
				backup_s3_bucket: "",
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

		updateEvent() {
			this.form.button_disabled = true;

			this.$store
				.dispatch("settings/updateSettingsAction", {
					s3Key: this.form.backup_s3_key,
					s3Secret: this.form.backup_s3_secret,
					s3Endpoint: this.form.backup_s3_endpoint,
					s3Region: this.form.backup_s3_region,
					s3Bucket: this.form.backup_s3_bucket,
				})
				.then(
					() => {
						this.$buefy.toast.open({
							message: "Settings updated successfully",
							type: "is-success",
						});

						this.form.button_disabled = false;
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

		this.$store.dispatch("settings/getSettingsAction").then(
			() => {
				let settings = this.$store.getters["settings/getSettingsResult"];

				this.form.backup_s3_key = settings.s3Key;
				this.form.backup_s3_secret = settings.s3Secret;
				this.form.backup_s3_endpoint = settings.s3Endpoint;
				this.form.backup_s3_region = settings.s3Region;
				this.form.backup_s3_bucket = settings.s3Bucket;

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

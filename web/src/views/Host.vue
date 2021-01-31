<!-- @format -->

<template>
	<div class="columns is-desktop is-centered">
		<div class="column is-9">
			<section>
				<b-table
					:data="crons.data"
					:paginated="crons.isPaginated"
					:per-page="crons.perPage"
					:current-page.sync="crons.currentPage"
					:pagination-simple="crons.isPaginationSimple"
					:pagination-position="crons.paginationPosition"
					:pagination-rounded="crons.isPaginationRounded"
					default-sort="name"
					aria-next-label="Next page"
					aria-previous-label="Previous page"
					aria-page-label="Page"
					aria-current-label="Current page"
				>
					<b-table-column
						field="name"
						label="Backup Name"
						centered
						v-slot="props"
					>
						<span class="tag is-light">
							{{ props.row.name }}
						</span>
					</b-table-column>

					<b-table-column field="jobs" label="Jobs" centered v-slot="props">
						<span class="tag is-warning is-light">
							Pending: {{ props.row.pendingJobs }}
						</span>
						-
						<span class="tag is-success is-light">
							Success: {{ props.row.successJobs }}
						</span>
						-
						<span class="tag is-danger is-light">
							Failed: {{ props.row.failedJobs }}
						</span>
					</b-table-column>

					<b-table-column
						field="status"
						label="Last Run"
						centered
						v-slot="props"
					>
						<span class="tag is-success is-light">
							{{ new Date(props.row.lastRun).toLocaleString() }}
						</span>
					</b-table-column>

					<b-table-column
						field="createdAt"
						label="Creation Date"
						centered
						v-slot="props"
					>
						<span class="tag is-info is-light">
							{{ new Date(props.row.createdAt).toLocaleDateString() }}
						</span>
					</b-table-column>

					<b-table-column label="Actions" centered v-slot="props">
						<b-button
							size="is-small"
							@click="editHostCronAction($route.params.hostId, props.row.id)"
							type="is-link is-warning is-light"
							>Edit</b-button
						>
						-
						<b-button
							size="is-small"
							type="is-link is-danger is-light"
							@click="deleteHostCronAction($route.params.hostId, props.row.id)"
							>Delete</b-button
						>
					</b-table-column>

					<td slot="empty" colspan="5">No records found.</td>
				</b-table>
			</section>

			<hr />
			<b-button
				type="is-success is-light"
				size="is-small"
				@click="newHostCronAction($route.params.hostId)"
				>New Backup Cron</b-button
			>
			-
			<b-button
				type="is-danger is-light"
				size="is-small"
				@click="deleteHostAction($route.params.hostId)"
				>Delete Host</b-button
			>

			<b-modal :active.sync="form.isActive" has-modal-card>
				<div class="modal-card" style="width: auto">
					<header class="modal-card-head">
						<p class="modal-card-title">{{ form.title }}</p>
					</header>
					<section class="modal-card-body">
						<b-field label="Backup Name">
							<b-input
								type="text"
								v-model="form.name"
								placeholder="Application Database"
								required
							>
							</b-input>
						</b-field>
						<b-field label="Backup Type">
							<b-select
								v-model="form.type"
								placeholder="Select a type"
								expanded
							>
								<option value="@BackupDirectory">Directory</option>
							</b-select>
						</b-field>
						<b-field label="Directory Path">
							<b-input
								type="text"
								v-model="form.directory"
								placeholder="/etc/backups/app_database"
								required
							>
							</b-input>
						</b-field>
						<b-field label="Backup Interval Type">
							<b-select
								v-model="form.intervalType"
								placeholder="Select a type"
								expanded
							>
								<option value="@second">Second</option>
								<option value="@minute">Minute</option>
								<option value="@hour">Hour</option>
								<option value="@day">Day</option>
								<option value="@month">Month</option>
							</b-select>
						</b-field>
						<b-field label="Backup Interval">
							<b-input
								type="number"
								v-model="form.interval"
								placeholder="30"
								required
							>
							</b-input>
						</b-field>
						<b-field label="Retention Days">
							<b-input
								type="number"
								v-model="form.retention"
								placeholder="10"
								required
							>
							</b-input>
						</b-field>
					</section>
					<footer class="modal-card-foot">
						<button class="button" type="button" @click="closeFormAction()">
							Close
						</button>
						<button class="button is-primary" @click="submitForm">
							Submit
						</button>
					</footer>
				</div>
			</b-modal>
		</div>
	</div>
</template>

<script>
export default {
	name: "host",

	data() {
		return {
			crons: {
				data: [],
				isPaginated: true,
				isPaginationSimple: false,
				isPaginationRounded: false,
				paginationPosition: "bottom",
				currentPage: 1,
				perPage: 15,
			},
			form: {
				isActive: false,
				title: "",
				name: "",
				interval: 30,
				retention: 10,
				directory: "",
				intervalType: "@minute",
				type: "@BackupDirectory",
				hostId: "",
				cronId: "",
			},
		};
	},
	methods: {
		deleteHostAction(hostId) {
			this.$buefy.dialog.confirm({
				message: "Are you sure?",
				onConfirm: () => {
					this.$store
						.dispatch("host/deleteHostAction", {
							hostname: hostId,
						})
						.then(
							() => {
								this.$buefy.toast.open({
									message: "Host deleted successfully",
									type: "is-success",
								});

								this.$router.push("/hosts");
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
							}
						);
				},
			});
		},

		deleteHostCronAction(hostId, cronId) {
			this.$buefy.dialog.confirm({
				message: "Are you sure?",
				onConfirm: () => {
					this.$store
						.dispatch("host/deleteHostCronAction", {
							hostname: hostId,
							cronId: cronId,
						})
						.then(
							() => {
								this.$buefy.toast.open({
									message: "Host cron deleted successfully",
									type: "is-success",
								});
								this.loadInitialState();
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
							}
						);
				},
			});
		},

		editHostCronAction(hostId, cronId) {
			this.form.title = "Edit Host Backup Cron";
			this.form.isActive = true;
			this.form.name = "";
			this.form.interval = 30;
			this.form.retention = 10;
			this.form.directory = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = hostId;
			this.form.cronId = cronId;

			this.$store
				.dispatch("host/getHostCronAction", {
					hostname: hostId,
					cronId: cronId,
				})
				.then(
					() => {
						let data = this.$store.getters["host/getHostCronResult"];

						if (data) {
							this.form.name = data.name;
							this.form.interval = data.interval;
							this.form.retention = data.request.retentionDays;
							this.form.directory = data.request.directory;
							this.form.intervalType = data.intervalType;
							this.form.type = data.request.type;
						}
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
					}
				);
		},

		newHostCronAction(hostId) {
			this.form.title = "New Host Backup Cron";
			this.form.isActive = true;
			this.form.name = "";
			this.form.interval = 30;
			this.form.retention = 10;
			this.form.directory = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = hostId;
			this.form.cronId = "";
		},

		closeFormAction() {
			this.form.isActive = false;
			this.form.title = "";
			this.form.name = "";
			this.form.interval = 30;
			this.form.retention = 10;
			this.form.directory = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = "";
			this.form.cronId = "";
		},

		submitForm() {
			let action = "host/createHostCronAction";

			if (this.form.cronId != "") {
				action = "host/updateHostCronAction";
			}

			this.$store
				.dispatch(action, {
					hostname: this.form.hostId,
					cronId: this.form.cronId,
					name: this.form.name,
					interval: this.form.interval.toString(),
					retention: this.form.retention.toString(),
					directory: this.form.directory,
					intervalType: this.form.intervalType,
					type: this.form.type,
				})
				.then(
					() => {
						this.$buefy.toast.open({
							message: "Backup cron submitted successfully",
							type: "is-success",
						});

						this.form.isActive = false;
						this.loadInitialState();
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
						this.form.isActive = false;
					}
				);
		},

		loadInitialState() {
			this.$store
				.dispatch("host/getHostCronsAction", {
					hostname: this.$route.params.hostId,
				})
				.then(
					() => {
						let data = this.$store.getters["host/getHostCronsResult"];

						if (data.crons) {
							this.crons.data = data.crons;
						} else {
							this.crons.data = [];
						}
					},
					(err) => {
						this.$buefy.toast.open({
							message: err.response.data.errorMessage,
							type: "is-danger is-light",
						});
					}
				);
		},
	},

	mounted() {
		this.loadInitialState();
	},
};
</script>

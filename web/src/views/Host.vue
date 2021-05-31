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
						field="id"
						label="Backup Cron ID"
						centered
						v-slot="props"
					>
						<span class="tag is-light">
							{{ props.row.id }}
						</span>
					</b-table-column>

					<b-table-column field="name" label="Name" centered v-slot="props">
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

					<td slot="empty" colspan="6">No records found.</td>
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
								<option value="@BackupMySQL">MySQL</option>
								<option value="@BackupSQLite">SQLite</option>
							</b-select>
						</b-field>

						<b-field label="Run Script Before">
							<b-input
								type="text"
								v-model="form.beforeScript"
								placeholder="/usr/bin/innobackupex --incremental"
							>
							</b-input>
						</b-field>

						<template v-if="form.type == '@BackupDirectory'">
							<b-field label="Directory Path">
								<b-input
									type="text"
									v-model="form.directory"
									placeholder="/etc/backups/app_database"
									required
								>
								</b-input>
							</b-field>
						</template>

						<template v-if="form.type == '@BackupSQLite'">
							<b-field label="SQLite Path">
								<b-input
									type="text"
									v-model="form.sqlitePath"
									placeholder="/etc/apps/db.sqlite3"
									required
								>
								</b-input>
							</b-field>
						</template>

						<template v-if="form.type == '@BackupMySQL'">
							<b-field label="MySQL Host">
								<b-input
									type="text"
									v-model="form.mysqlHost"
									placeholder="127.0.0.1"
									required
								>
								</b-input>
							</b-field>
							<b-field label="MySQL Port">
								<b-input
									type="text"
									v-model="form.mysqlPort"
									placeholder="3306"
									required
								>
								</b-input>
							</b-field>
							<b-field label="MySQL Username">
								<b-input
									type="text"
									v-model="form.mysqlUsername"
									placeholder="root"
									required
								>
								</b-input>
							</b-field>
							<b-field label="MySQL Password">
								<b-input
									type="text"
									v-model="form.mysqlPassword"
									placeholder="root"
									required
								>
								</b-input>
							</b-field>

							<b-field label="MySQL Database Name">
								<b-input
									type="text"
									v-model="form.mysqlDatabase"
									placeholder=""
									required
								>
								</b-input>
							</b-field>
							<b-field label="MySQL Table Name">
								<b-input
									type="text"
									v-model="form.mysqlTable"
									placeholder=""
									required
								>
								</b-input>
							</b-field>
							<b-field label="Backup All Databases">
								<b-select
									v-model="form.mysqlAllDatabases"
									placeholder="Select to backup all"
									expanded
								>
									<option value="false">No</option>
									<option value="true">Yes</option>
								</b-select>
							</b-field>
							<b-field label="MySQL Dump Options">
								<b-input
									type="text"
									v-model="form.mysqlOptions"
									placeholder="--single-transaction,--quick,--lock-tables=false"
									required
								>
								</b-input>
							</b-field>
						</template>

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
				isPaginationRounded: true,
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
				beforeScript: "",
				directory: "",

				mysqlHost: "127.0.0.1",
				mysqlPort: "3306",
				mysqlUsername: "root",
				mysqlPassword: "root",
				mysqlAllDatabases: "false",
				mysqlDatabase: "",
				mysqlTable: "",
				mysqlOptions: "--single-transaction,--quick,--lock-tables=false",

				intervalType: "@minute",
				type: "@BackupDirectory",
				hostId: "",
				cronId: "",

				sqlitePath: "/etc/apps/db.sqlite3",
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
			this.form.beforeScript = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = hostId;
			this.form.cronId = cronId;
			this.form.sqlitePath = "/etc/apps/db.sqlite3";

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
							this.form.beforeScript = data.request.beforeScript;
							this.form.directory = data.request.directory;

							this.form.mysqlHost = data.request.mysqlHost;
							this.form.mysqlPort = data.request.mysqlPort;
							this.form.mysqlUsername = data.request.mysqlUsername;
							this.form.mysqlPassword = data.request.mysqlPassword;
							this.form.mysqlAllDatabases = data.request.mysqlAllDatabases;
							this.form.mysqlDatabase = data.request.mysqlDatabase;
							this.form.mysqlTable = data.request.mysqlTable;
							this.form.mysqlOptions = data.request.mysqlOptions;
							this.form.sqlitePath = data.request.sqlitePath;

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
			this.form.beforeScript = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = hostId;
			this.form.cronId = "";

			this.form.mysqlHost = "127.0.0.1";
			this.form.mysqlPort = "3306";
			this.form.mysqlUsername = "root";
			this.form.mysqlPassword = "root";
			this.form.mysqlAllDatabases = "false";
			this.form.mysqlDatabase = "";
			this.form.mysqlTable = "";
			this.form.mysqlOptions =
				"--single-transaction,--quick,--lock-tables=false";

			this.form.sqlitePath = "/etc/apps/db.sqlite3";
		},

		closeFormAction() {
			this.form.isActive = false;
			this.form.title = "";
			this.form.name = "";
			this.form.interval = 30;
			this.form.retention = 10;
			this.form.directory = "";
			this.form.beforeScript = "";
			this.form.intervalType = "@minute";
			this.form.type = "@BackupDirectory";
			this.form.hostId = "";
			this.form.cronId = "";

			this.form.mysqlHost = "127.0.0.1";
			this.form.mysqlPort = "3306";
			this.form.mysqlUsername = "root";
			this.form.mysqlPassword = "root";
			this.form.mysqlAllDatabases = "false";
			this.form.mysqlDatabase = "";
			this.form.mysqlTable = "";
			this.form.mysqlOptions =
				"--single-transaction,--quick,--lock-tables=false";

			this.form.sqlitePath = "/etc/apps/db.sqlite3";
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
					beforeScript: this.form.beforeScript,
					directory: this.form.directory,
					intervalType: this.form.intervalType,

					mysqlHost: this.form.mysqlHost,
					mysqlPort: this.form.mysqlPort,
					mysqlUsername: this.form.mysqlUsername,
					mysqlPassword: this.form.mysqlPassword,
					mysqlAllDatabases: this.form.mysqlAllDatabases,
					mysqlDatabase: this.form.mysqlDatabase,
					mysqlTable: this.form.mysqlTable,
					mysqlOptions: this.form.mysqlOptions,

					sqlitePath: this.form.sqlitePath,

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
			this.loading();

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
	},

	mounted() {
		this.loadInitialState();
	},
};
</script>

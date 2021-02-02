<!-- @format -->

<template>
	<div class="columns is-desktop is-centered">
		<div class="column is-9">
			<section>
				<b-table
					:data="data"
					:paginated="isPaginated"
					:per-page="perPage"
					:current-page.sync="currentPage"
					:pagination-simple="isPaginationSimple"
					:pagination-position="paginationPosition"
					:pagination-rounded="isPaginationRounded"
					default-sort="hostname"
					aria-next-label="Next page"
					aria-previous-label="Previous page"
					aria-page-label="Page"
					aria-current-label="Current page"
				>
					<b-table-column field="hostname" label="Host" centered v-slot="props">
						<span class="tag is-light">
							{{ props.row.hostname }}
						</span>
					</b-table-column>

					<b-table-column
						field="onlineAgents"
						label="Online Agents"
						centered
						v-slot="props"
					>
						<span class="tag is-success is-light">
							{{ props.row.onlineAgents }}
						</span>
					</b-table-column>

					<b-table-column
						field="createdAt"
						label="Creation Date"
						centered
						v-slot="props"
					>
						<span class="tag is-warning is-light">
							{{ new Date(props.row.createdAt).toLocaleString() }}
						</span>
					</b-table-column>

					<b-table-column label="Actions" centered v-slot="props">
						<b-button
							tag="router-link"
							size="is-small"
							:to="'/hosts/' + props.row.hostname"
							type="is-link is-danger is-light"
							>View</b-button
						>
					</b-table-column>
					<td slot="empty" colspan="4">No records found.</td>
				</b-table>
			</section>
		</div>
	</div>
</template>

<script>
export default {
	name: "hosts",

	data() {
		return {
			data: [],
			isPaginated: true,
			isPaginationSimple: false,
			isPaginationRounded: true,
			paginationPosition: "bottom",
			currentPage: 1,
			perPage: 15,

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
				container: this.loader.isFullPage ? null : this.$refs.element.$el
			});
		},
	},

	mounted() {
		this.loading();

		this.$store.dispatch("host/getHostsAction").then(
			() => {
				let data = this.$store.getters["host/getHostsResult"].hosts;

				if (data) {
					this.data = data;
				} else {
					this.data = [];
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

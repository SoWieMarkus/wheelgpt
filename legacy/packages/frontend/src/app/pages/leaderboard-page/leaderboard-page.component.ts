import { Component, inject, model, signal } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTableModule } from "@angular/material/table";
import { Title } from "@angular/platform-browser";
import { ActivatedRoute, Router } from "@angular/router";
import { TranslatePipe, TranslateService } from "@ngx-translate/core";
import { LoadingComponent } from "../../components/loading/loading.component";
import type { PublicChannelDetails } from "../../schemas/channel";
import type { LeaderboardEntry } from "../../schemas/leaderboard";
import { BackendService } from "../../services/backend.service";

@Component({
	selector: "app-leaderboard-page",
	imports: [
		LoadingComponent,
		MatButtonModule,
		FormsModule,
		MatIconModule,
		MatTableModule,
		MatInputModule,
		MatFormFieldModule,
		TranslatePipe,
	],
	templateUrl: "./leaderboard-page.component.html",
	styleUrl: "./leaderboard-page.component.scss",
})
export class LeaderboardPage {
	private readonly backendService = inject(BackendService);
	private readonly snackbar = inject(MatSnackBar);
	private readonly route = inject(ActivatedRoute);
	private readonly router = inject(Router);

	protected readonly channel = signal<PublicChannelDetails | null>(null);
	protected readonly leaderboard = signal<LeaderboardEntry[]>([]);
	protected readonly channelId = signal<string | null>(null);

	// Tracks the current page of the leaderboard
	protected readonly currentPage = signal<number>(1);

	// This signal indicates whether there are more pages to load
	// It is set to false when the last page that was loaded has no entries
	protected readonly hasNext = signal<boolean>(true);

	// This signal indicates whether the leaderboard is currently loading
	protected readonly isLeaderboardLoading = signal<boolean>(true);

	protected readonly searchUsername = model("");
	protected readonly isLoadingUser = signal<boolean>(false);
	protected readonly userResult = signal<LeaderboardEntry | null>(null);

	private readonly translateService = inject(TranslateService);
	private readonly title = inject(Title);

	public ngOnInit(): void {
		this.route.params.subscribe({
			next: (params) => {
				window.scrollTo(0, 0);
				// biome-ignore lint/complexity/useLiteralKeys: Doesn't work here
				const channelId = params["channelId"];
				if (!channelId) {
					this.router.navigate(["/"]);
					return;
				}
				this.channelId.set(channelId);
				this.currentPage.set(1);
				this.leaderboard.set([]);
				this.channel.set(null);
				this.hasNext.set(true);
				this.isLeaderboardLoading.set(true);
				this.searchUsername.set("");

				this.backendService.channel
					.getChannelById(channelId)
					.then((data) => {
						this.title.setTitle(`WheelGPT: ${data.displayName}`);
						this.channel.set(data);
					})
					.catch((error) => {
						console.error("Failed to load channel details:", error);
						this.channel.set(null);
						const message = this.translateService.instant("pages.leaderboard.actions.load-channel.error");
						this.snackbar.open(message, "", { duration: 3000 });
						this.router.navigate(["/"]);
					});
				this.backendService.leaderboard
					.get(channelId)
					.then((data) => {
						this.leaderboard.set(data);
					})
					.catch((error) => {
						console.error("Failed to load leaderboard:", error);
						const message = this.translateService.instant("pages.leaderboard.actions.load-leaderboard.error");
						this.snackbar.open(message, "", { duration: 3000 });
						this.router.navigate(["/"]);
					})
					.finally(() => {
						this.isLeaderboardLoading.set(false);
					});
				this.searchUserDetails();
			},
			error: (error) => {
				console.error("Error fetching route parameters:", error);
				this.router.navigate(["/"]);
			},
		});
	}

	public loadMore(): void {
		// Prevent loading more if already loading or no more pages available
		if (this.isLeaderboardLoading() || !this.hasNext()) {
			return;
		}

		this.isLeaderboardLoading.set(true);
		const currentPage = this.currentPage();
		const nextPage = currentPage + 1;
		this.currentPage.set(nextPage);

		const channelId = this.channelId();
		if (!channelId) {
			throw new Error("Channel ID is not set");
		}
		this.backendService.leaderboard
			.get(channelId, nextPage)
			.then((data) => {
				if (data.length === 0) {
					this.hasNext.set(false);
					return;
				}
				const currentLeaderboard = this.leaderboard();
				this.leaderboard.set([...currentLeaderboard, ...data]);
			})
			.catch((error) => {
				console.error("Failed to load more leaderboard entries:", error);
				const message = this.translateService.instant("pages.leaderboard.actions.load-leaderboard.error");
				this.snackbar.open(message, "", { duration: 3000 });
			})
			.finally(() => {
				this.isLeaderboardLoading.set(false);
			});
	}

	// When a user enters a username in the search field, this method is called
	// It searches for the user in the leaderboard and displays their details if found
	public searchUserDetails(): void {
		const channelId = this.channelId();
		if (!channelId) {
			throw new Error("Channel ID is not set");
		}
		const username = this.searchUsername();
		if (!username) {
			this.userResult.set(null);
			return;
		}

		this.isLoadingUser.set(true);
		this.backendService.leaderboard
			.getByName(channelId, username)
			.then((data) => {
				this.userResult.set(data);
			})
			.catch((error) => {
				console.error("Failed to search user:", error);
				this.userResult.set(null);
				const message = this.translateService.instant("pages.leaderboard.actions.search.error");
				this.snackbar.open(message, "", { duration: 3000 });
			})
			.finally(() => {
				this.isLoadingUser.set(false);
			});
	}
}

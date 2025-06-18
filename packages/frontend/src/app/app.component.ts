import { Component, inject } from "@angular/core";
import { Router, RouterOutlet } from "@angular/router";
import { AuthenticationService } from "./services/authentication.service";
import { DomSanitizer } from "@angular/platform-browser";
import { MatIconModule, MatIconRegistry } from '@angular/material/icon';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { MatDividerModule } from "@angular/material/divider";
import { ProfileService } from "./services/profile.service";

@Component({
  selector: "app-root",
  imports: [RouterOutlet, MatDividerModule, TranslatePipe, MatIconModule],
  templateUrl: "./app.component.html",
  styleUrl: "./app.component.scss",
})
export class AppComponent {
  protected readonly authenticationService = inject(AuthenticationService);
  protected readonly router = inject(Router);
  private readonly sanitizer = inject(DomSanitizer);
  private readonly registry = inject(MatIconRegistry);
  private readonly translateService = inject(TranslateService);
  protected readonly year = new Date().getFullYear();
  protected readonly profileService = inject(ProfileService);

  constructor() {
    this.initializeLanguageService();
    this.initializeIcons();
  }

  private initializeLanguageService() {
    this.translateService.setDefaultLang("en");
    const cachedLanguage = localStorage.getItem('language');
    if (cachedLanguage) {
      this.translateService.use(cachedLanguage);
      return;
    }
    const browserLanguage = this.translateService.getBrowserLang();
    // Use browser language if available, otherwise default to English
    this.translateService.use(browserLanguage?.match(/de|en/) ? browserLanguage : 'en');
  }

  public switchLanguage(lang: "de" | "en") {
    this.translateService.use(lang);
  }

  private initializeIcons() {
    const icons = ["github", "twitch"];
    for (const icon of icons) {
      this.registry.addSvgIcon(icon, this.sanitizer.bypassSecurityTrustResourceUrl(`icons/${icon}.svg`));
    }
  }

  protected logout(): void {
    this.authenticationService.removeToken();
    this.router.navigate(["/landing"]);
  }

  public ngOnInit(): void {
    if (this.authenticationService.isTokenExpired()) {
      this.authenticationService.removeToken();
      return;
    }
    this.authenticationService.initialize();
  }

}

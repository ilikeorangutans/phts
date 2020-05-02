import { Component } from '@angular/core';
import { AuthService } from './services/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  authenticated: boolean = false;

  constructor(private readonly authService: AuthService) {
    this.authService.authenticated.subscribe((authenticated) => {
      this.authenticated = authenticated;
    });
  }

  login() {
    this.authService.authenticate('username', 'password');
  }

  logout() {
    this.authService.logout();
  }
}

import { Component } from '@angular/core';
import { SessionService } from './services/session.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  authenticated: boolean = false;

  constructor(readonly sessionService: SessionService) {
    this.sessionService.hasSession.subscribe((hasSession) => {
      this.authenticated = hasSession;
    });
  }

  login() {
    this.sessionService.start('test@test.local', 'test');
  }

  logout() {
    this.sessionService.destroy();
  }
}

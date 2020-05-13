import { Component } from '@angular/core';
import { SessionService } from './services/session.service';
import { HttpClient } from '@angular/common/http';
import { PathService } from './services/path.service';
import { Observable } from 'rxjs';
import { Store, select } from '@ngrx/store';
import { SessionState } from './reducers/session.reducer';
import { login, logout } from './actions/session.action';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  authenticated: boolean = false;

  constructor(
    readonly sessionService: SessionService,
    readonly http: HttpClient,
    readonly pathService: PathService,
    private store: Store<{ session: SessionState }>
  ) {
    this.sessionService.hasSession.subscribe((hasSession) => {
      this.authenticated = hasSession;
    });
    this.sessionState$ = store.pipe(select('session'))
  }

  login() {
    this.sessionService.start('test@test.local', 'test');
  }

  logout() {
    this.sessionService.destroy();
  }

  sessionState$: Observable<SessionState>;
  count$: Observable<number>;

  xxxlogin() {
    this.store.dispatch(login({ username: 'user', password: 'password' }));
  }

  xxxlogout() {
    this.store.dispatch(logout());
  }

}

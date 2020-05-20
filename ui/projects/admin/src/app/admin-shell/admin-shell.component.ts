import { Component, OnInit } from '@angular/core';
import { SessionService } from '../services/session.service';
import { Router } from '@angular/router';
import { PhtsService } from '../services/phts.service';
import { fromEvent } from 'rxjs';
import { map, distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-admin-shell',
  templateUrl: './admin-shell.component.html',
  styleUrls: ['./admin-shell.component.css'],
})
export class AdminShellComponent implements OnInit {
  navBarClasses = {
    show: false,
  };
  navItems: Array<NavItem> = [];
  navCollapsed = false;

  constructor(
    private sessionService: SessionService,
    private router: Router,
    readonly phtsService: PhtsService
  ) {
    fromEvent(window, 'resize')
      .pipe(
        map((_) => window.innerWidth),
        map((width) => width < 576),
        distinctUntilChanged()
      )
      .subscribe((isSmall) => {
        this.navBarClasses = {
          show: false,
        };
        this.navCollapsed = isSmall;
      });
  }

  ngOnInit() {
    this.navItems = [
      new NavItem('Photos', 'collection'),
      new NavItem('Share Sites', 'share-site'),
      new NavItem('Account', 'account'),
    ];

    this.phtsService.refreshVersion();
  }

  logout() {
    this.sessionService.logout();
    this.router.navigate(['login']);
  }

  toggleNav(): void {
    this.navBarClasses['show'] = !this.navBarClasses['show'];
  }
}

export class NavItem {
  constructor(public readonly title: string, public readonly link: string) {}
}

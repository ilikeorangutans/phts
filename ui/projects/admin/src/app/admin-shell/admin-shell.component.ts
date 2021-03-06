import { Component, OnInit } from '@angular/core';
import { SessionService } from '../services/session.service';
import { Router } from '@angular/router';
import { PhtsService } from '../services/phts.service';

@Component({
  selector: 'app-admin-shell',
  templateUrl: './admin-shell.component.html',
  styleUrls: ['./admin-shell.component.css'],
})
export class AdminShellComponent implements OnInit {
  navBarClasses = {
    'is-active': false,
  };
  navItems: Array<NavItem> = [];
  navCollapsed = false;

  constructor(
    readonly sessionService: SessionService,
    private router: Router,
    readonly phtsService: PhtsService
  ) { }

  ngOnInit() {
    this.navItems = [
      new NavItem('Collections', 'collection'),
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
    this.navBarClasses['is-active'] = !this.navBarClasses['is-active'];
  }
}

export class NavItem {
  constructor(public readonly title: string, public readonly link: string) { }
}

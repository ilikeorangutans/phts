import { PhtsService } from './../services/phts.service';
import { Observable } from 'rxjs/Observable';
import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { SessionService } from '../services/session.service';

@Component({
  selector: 'admin-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  navItems: Array<NavItem> = [];

  constructor(
    private sessionService: SessionService,
    private router: Router,
    readonly phtsService: PhtsService
  ) { }

  ngOnInit() {
    this.navItems = [
      new NavItem('Photos', 'collection'),
      new NavItem('Share Sites', 'share-site'),
      new NavItem('Account', 'account')
    ];

    this.phtsService.refreshVersion();
  }

  logout() {
    this.sessionService.logout();
    this.router.navigate(['admin']);
  }
}

export class NavItem {

  constructor(
    public title: string,
    public link: string
  ) { }
}

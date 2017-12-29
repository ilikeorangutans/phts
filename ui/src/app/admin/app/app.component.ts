import { Router } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { SessionService } from '../services/session.service';

@Component({
  selector: 'admin-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  constructor(
    private sessionService: SessionService,
    private router: Router
  ) { }

  ngOnInit() {
  }


  logout() {
    this.sessionService.logout();
    this.router.navigate(['admin']);
  }
}

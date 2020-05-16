import { Component, OnInit } from '@angular/core';
import { SessionService } from '../services/session.service';

class Credentials {
  email: string;
  password: string;
}

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  credentials: Credentials = new Credentials();

  constructor(
    private readonly sessionService: SessionService
  ) { }

  ngOnInit(): void {
  }

  onSubmit() {
    this.sessionService.start(this.credentials.email, this.credentials.password);
  }
}

import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';
import { PathService } from '../services/path.service';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

@Component({
  selector: 'app-join',
  templateUrl: './join.component.html',
  styleUrls: ['./join.component.css'],
})
export class JoinComponent implements OnInit {
  invitation$: Observable<Invitation>;

  joinRequest: JoinRequest = new JoinRequest();

  submitting = false;

  constructor(
    private paths: PathService,
    private http: HttpClient,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((x) => {
      console.log(x);
      const url = this.paths.joinToken(x.token);
      console.log(url);

      this.invitation$ = this.http.get<Invitation>(url).pipe(
        tap((invitation) => {
          this.joinRequest.email = invitation.email;
          this.joinRequest.token = x.token;
        })
      );
    });
  }

  onSubmit() {
    this.submitting = true;
    const url = this.paths.joinToken(this.joinRequest.token);
    console.log(url);

    this.http.post<JoinResponse>(url, this.joinRequest).subscribe((resp) => {
      console.log(resp);
      this.submitting = false;
    });
  }
}

class Invitation {
  email: string;
}

class JoinRequest {
  email: string;
  name: string;
  password: string;
  confirmPassword: string;
  token: string;
}

class JoinResponse {}

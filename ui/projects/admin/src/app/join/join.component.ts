import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { PathService } from '../services/path.service';
import { Observable, NEVER } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

@Component({
  selector: 'app-join',
  templateUrl: './join.component.html',
  styleUrls: ['./join.component.css'],
})
export class JoinComponent implements OnInit {
  invitation$: Observable<Invitation>;

  loading = true;
  submitting = false;
  error = '';

  constructor(
    private paths: PathService,
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((x) => {
      const url = this.paths.joinToken(x.token);
      this.invitation$ = this.http.get<Invitation>(url).pipe(
        retry(2),
        catchError((err) => this.handleError(err))
      );
    });
  }

  onSubmit(joinRequest: JoinRequest) {
    this.submitting = true;
    const url = this.paths.joinToken(joinRequest.token);

    this.http.post<JoinResponse>(url, joinRequest).subscribe((resp) => {
      console.log(resp);
      this.submitting = false;
      this.router.navigate(['/']);
    });
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      console.log('error event');
    } else {
      console.log('error code', error.error, error.status, error.statusText);
    }

    this.loading = false;
    this.error = 'Invitation does not exist.';
    return NEVER;
  }
}

export class Invitation {
  email: string;
  token: string;
}

export class JoinRequest {
  email: string;
  name: string;
  password: string;
  confirmPassword: string;
  token: string;
}

class JoinResponse {}

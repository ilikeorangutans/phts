import { HttpClient } from '@angular/common/http';
import { PathService } from './../../services/path.service';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable()
export class AccountService {
  constructor(private http: HttpClient, private pathService: PathService) {}

  updatePassword(
    existingPassword: string,
    newPassword: string
  ): Observable<null> {
    const url = this.pathService.changePassword();
    const change = new PasswordChange(newPassword, existingPassword);

    return this.http.post(url, change).pipe(map((_) => null));
  }
}

class PasswordChange {
  password: string;
  oldPassword: string;

  constructor(p: string, op: string) {
    this.password = p;
    this.oldPassword = op;
  }
}

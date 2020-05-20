import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { PhtsVersion } from './../models/phts-version';
import { PathService } from './path.service';
import { Subject, BehaviorSubject, Observable } from 'rxjs';
import { filter } from 'rxjs/operators';

@Injectable()
export class PhtsService {
  private readonly _version: Subject<PhtsVersion | null> = new BehaviorSubject<PhtsVersion | null>(
    null
  );

  readonly version: Observable<PhtsVersion | null> = this._version
    .asObservable()
    .pipe(filter((v) => v !== null));

  constructor(private pathService: PathService, private http: HttpClient) {}

  refreshVersion(): void {
    const url = this.pathService.version();
    this.http
      .get<PhtsVersion>(url)
      .subscribe((version) => this._version.next(version));
  }
}

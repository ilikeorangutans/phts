import { ShareSiteService } from './../services/share-site.service';
import { ShareSite } from '../models/share-site';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { first } from 'rxjs/operators';

@Injectable()
export class ShareSiteStore {
  private readonly _all = new BehaviorSubject<Array<ShareSite>>([]);

  readonly all: Observable<Array<ShareSite>> = this._all.asObservable();

  constructor(private readonly shareSiteService: ShareSiteService) {}

  refresh(): void {
    this.shareSiteService
      .list()
      .pipe(first())
      .subscribe((sites) => this._all.next(sites));
  }

  delete(shareSite: ShareSite): void {
    console.log('implement me! delete', shareSite);
  }

  save(shareSite: ShareSite): void {
    console.log('saving share site');
    this.shareSiteService
      .save(shareSite)
      .pipe(first())
      .subscribe((_) => this.refresh());
  }
}

import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import 'rxjs/add/operator/switchMap';

import { ShareService } from './../services/share.service';
import { Share } from './../models/share';

@Component({
  selector: 'public-share-viewer',
  templateUrl: './share-viewer.component.html',
  styleUrls: ['./share-viewer.component.css']
})
export class ShareViewerComponent implements OnInit, OnDestroy {

  private sub: Subscription;

  constructor(
    private shareService: ShareService,
    private route: ActivatedRoute,
    private router: Router
  ) { }

  share: Share;

  ngOnInit() {

    this.sub = this.route.params
      .map(params => params['slug'] as string)
      .switchMap(slug => this.shareService.forSlug(slug))
      .subscribe(share => this.share = share);
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }
}

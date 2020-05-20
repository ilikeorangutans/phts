import { Component, OnInit } from '@angular/core';
import { ShareService } from '../services/share.service';
import { ActivatedRoute } from '@angular/router';
import { map, switchMap } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { Share } from '../models/share';

@Component({
  selector: 'app-share-viewer',
  templateUrl: './share-viewer.component.html',
  styleUrls: ['./share-viewer.component.css'],
})
export class ShareViewerComponent implements OnInit {
  constructor(private shares: ShareService, private route: ActivatedRoute) {}

  share$: Observable<Share>;

  ngOnInit(): void {
    this.share$ = this.route.paramMap.pipe(
      map((params) => params.get('slug') as string),
      switchMap((slug) => this.shares.forSlug(slug))
    );
    // TODO here we should get the index of the current photo from the param map
  }
}

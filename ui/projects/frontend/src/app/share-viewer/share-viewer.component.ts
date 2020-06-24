import { Component, OnInit, HostListener } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { map, switchMap, tap } from 'rxjs/operators';
import { Observable } from 'rxjs';

import { ShareService } from '../services/share.service';
import { Share } from '../models/share';

@Component({
  selector: 'app-share-viewer',
  templateUrl: './share-viewer.component.html',
  styleUrls: ['./share-viewer.component.css'],
})
export class ShareViewerComponent implements OnInit {
  constructor(private shares: ShareService, private route: ActivatedRoute) { }

  screenWidth: number;
  screenHeight: number;
  photoIndex = 0;

  share$: Observable<Share>;

  renditionIndex = 0;

  ngOnInit(): void {
    this.share$ = this.route.paramMap.pipe(
      map((params) => params.get('slug') as string),
      switchMap((slug) => this.shares.forSlug(slug)),
      tap(share => this.findBestFit(share))
    );
  }

  findBestFit(share: Share) {
    const screenArea = window.innerHeight * window.innerWidth;
    const photo = share.photos[this.photoIndex];
    const bestFit = photo.renditions.sort((a, b) => Math.abs(screenArea - (a.width * a.height)) - Math.abs(screenArea - (b.width * b.height)))[0];
    this.renditionIndex = photo.renditions.findIndex((rendition, _) => rendition.id === bestFit?.id);
  }

  @HostListener('window:resize', ['$event'])
  onScreenSizeChange(_) {
    this.screenHeight = window.innerHeight;
    this.screenWidth = window.innerWidth;
  }

  switchRendition(index: number) {
    this.renditionIndex = index;
  }
}

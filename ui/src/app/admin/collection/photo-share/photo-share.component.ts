import { Component, Input, OnInit } from '@angular/core';

import { Collection } from './../../models/collection';
import { Photo } from './../../models/photo';
import { Share } from './../../models/share';
import { ShareService } from '../../services/share.service';
import { ShareSiteService } from './../../services/share-site.service';

@Component({
  selector: 'collection-photo-share',
  templateUrl: './photo-share.component.html',
  styleUrls: ['./photo-share.component.css']
})
export class PhotoShareComponent implements OnInit {

  @Input() collection: Collection;
  @Input() photo: Photo;
  share: Share = new Share();
  shareSites: Array<ShareSite> = [];

  constructor(
    private shareService: ShareService,
    private shareSiteService: ShareSiteService
  ) { }

  ngOnInit() {
    this.shareSiteService.list().then(sites => this.shareSites = sites);
  }

  onSubmit() {
    this.shareService.save(this.collection, this.photo, this.share);
    this.share = new Share();
  }
}

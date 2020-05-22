import { Component, Input, OnInit } from '@angular/core';

import { PhotoShares, ShareRequest } from './../../services/share.service';
import { Collection } from './../../models/collection';
import { Photo } from './../../models/photo';
import { ShareSiteService } from './../../services/share-site.service';
import { ShareSite } from '../../models/share-site';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { Observable } from 'rxjs';

@Component({
  selector: 'collection-photo-share',
  templateUrl: './photo-share.component.html',
  styleUrls: ['./photo-share.component.css'],
})
export class PhotoShareComponent implements OnInit {
  @Input() photoShares: PhotoShares;
  @Input() collection: Collection;
  @Input() photo: Photo;
  @Input() renditionConfigurations: Array<RenditionConfiguration> = [];
  share: ShareRequest;
  shareSites: Observable<Array<ShareSite>>;

  constructor(private shareSiteService: ShareSiteService) {}

  ngOnInit() {
    this.shareSites = this.shareSiteService.list();
    this.share = new ShareRequest(this.photo.id);
  }

  onSubmit() {
    const x = this.photoShares.add(this.share);
    x.subscribe((bla) => console.log('share added', bla));
    this.share = new ShareRequest(this.photo.id);
  }
}

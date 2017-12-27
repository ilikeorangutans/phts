import { ShareSite } from './../../models/share-site';
import { ShareSiteService } from './../../services/share-site.service';
import { PathService } from './../../services/path.service';
import { Rendition } from './../../models/rendition';
import { Photo } from './../../models/photo';
import { ShareService } from './../../services/share.service';
import { Component, OnInit, Input } from '@angular/core';
import { Share } from '../../models/share';
import { Collection } from '../../models/collection';

@Component({
  selector: 'app-share-photo',
  templateUrl: './share-photo.component.html',
  styleUrls: ['./share-photo.component.css']
})
export class SharePhotoComponent implements OnInit {

  @Input() photo: Photo;
  @Input() collection: Collection;

  share: Share = new Share();
  adminPreviewConfigID: number;
  shareSites: Array<ShareSite> = [];

  constructor(
    private pathService: PathService,
    private shareService: ShareService,
    private shareSiteService: ShareSiteService
  ) { }

  ngOnInit() {
    this.adminPreviewConfigID = this.collection.renditionConfigurations.find(rc => rc.name === 'admin thumbnails').id;
    this.shareSiteService.list().then(sites => this.shareSites = sites);
    this.share.photoID = this.photo.id;
  }

  onSubmit() {
    this.shareService.save(this.collection, this.photo, this.share);
  }

  adminPreview(photo: Photo): Rendition {
    return photo.renditions.find(r => r.renditionConfigurationID === this.adminPreviewConfigID);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

  debug(): String {
    return JSON.stringify(this.share);
  }
}

import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Component, Input, OnInit } from '@angular/core';

import { Photo } from './../../models/photo';
import { Collection } from '../../models/collection';
import { Share } from './../../models/share';
import { ShareService } from './../../services/share.service';

@Component({
  selector: 'collection-photo-share-list',
  templateUrl: './photo-share-list.component.html',
  styleUrls: ['./photo-share-list.component.css']
})
export class PhotoShareListComponent implements OnInit {

  @Input() collection: Collection;
  @Input() photo: Photo;

  shares: Array<Share> = [];

  constructor(
    private shareService: ShareService
  ) { }

  ngOnInit() {
    this.shareService
      .listForPhoto(this.collection, this.photo)
      .then(shares => this.shares = shares);
  }

  describeRendition(config: RenditionConfiguration): string {
    if (config.resize) {
      return `resize to ${config.width}×${config.height} at most`;
    }
    return config.name;
  }

  renditionClasses(config: RenditionConfiguration) {
    return {
      'badge-secondary': config.resize,
      'badge-primary': !config.resize
    };
  }

}

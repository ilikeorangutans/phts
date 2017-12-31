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
    // TODO is there a chance that collection/photo is not initialized here?
    this.shareService
      .listForPhoto(this.collection, this.photo)
      .then(shares => this.shares = shares);
  }

}

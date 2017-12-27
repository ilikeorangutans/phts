import { Collection } from './../../models/collection';
import { Share } from './../../models/share';
import { ShareService } from './../../services/share.service';
import { Photo } from './../../models/photo';
import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-photo-share-list',
  templateUrl: './photo-share-list.component.html',
  styleUrls: ['./photo-share-list.component.css']
})
export class PhotoShareListComponent implements OnInit {

  @Input() photo: Photo;
  @Input() collection: Collection;

  shares: Array<Share> = [];

  constructor(
    private shareService: ShareService
  ) { }

  ngOnInit() {
    this.shareService.listForPhoto(this.collection, this.photo)
      .then(result => this.shares = result);
  }

}

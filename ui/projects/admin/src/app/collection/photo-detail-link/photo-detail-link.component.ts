import { Component, OnInit, Input } from '@angular/core';

import { Album } from './../../models/album';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';

@Component({
  selector: 'photo-detail-link',
  templateUrl: './photo-detail-link.component.html',
  styleUrls: ['./photo-detail-link.component.css'],
})
export class PhotoDetailLinkComponent implements OnInit {
  @Input() photo: Photo;

  @Input() collection: Collection;

  @Input() album: Album;

  ngOnInit() {}

  detailURL(): string {
    if (this.album !== undefined) {
      return '';
    }

    return ['/collection', this.collection.slug, 'photos', this.photo.id].join(
      '/'
    );
  }
}

import { Rendition } from './../../models/rendition';
import { PathService } from './../../services/path.service';
import { Collection } from './../../models/collection';
import { Component, OnInit, Input } from '@angular/core';
import { Photo } from './../../models/photo';

@Component({
  selector: 'app-photo-thumbnail',
  templateUrl: './photo-thumbnail.component.html',
  styleUrls: ['./photo-thumbnail.component.css']
})
export class PhotoThumbnailComponent implements OnInit {

  @Input() photo: Photo;

  @Input() collection: Collection;

  private rendition: Rendition;

  constructor(
    private pathService: PathService
  ) { }

  ngOnInit() {
  }

  private previewID(): number {
    return this.collection.renditionConfigurations.find(config => config.name === 'admin thumbnails').id;
  }

  previewRendition(): Rendition {
    if (this.rendition !== undefined) {
      return this.rendition;
    }
    this.rendition = this.photo.renditions.find(rendition => rendition.renditionConfigurationID === this.previewID());

    return this.rendition;
  }

  thumbnailURL(): string {
    return this.pathService.rendition(this.collection, this.previewRendition());
  }

}

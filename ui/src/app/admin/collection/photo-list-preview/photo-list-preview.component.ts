import { Component, OnInit, Input } from '@angular/core';

import { Photo } from './../../models/photo';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { Rendition } from '../../models/rendition';
import { PathService } from '../../services/path.service';

@Component({
  selector: 'app-photo-list-preview',
  templateUrl: './photo-list-preview.component.html',
  styleUrls: ['./photo-list-preview.component.css']
})
export class PhotoListPreviewComponent implements OnInit {

  @Input() photo: Photo;

  @Input() renditionConfiguration: RenditionConfiguration;

  constructor(
    private pathService: PathService
  ) { }

  ngOnInit() {
  }

  viewURL(): string {
    return this.pathService.rendition(this.photo.collection, this.rendition());
  }

  rendition(): Rendition {
    return this.photo.renditions.find(r => r.renditionConfigurationID === this.renditionConfiguration.id);
  }

}

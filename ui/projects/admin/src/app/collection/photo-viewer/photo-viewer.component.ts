import { Component, OnInit, Input } from '@angular/core';

import { fromEvent } from 'rxjs';
import { map } from 'rxjs/operators';

import { Rendition } from './../../models/rendition';
import { PathService } from './../../services/path.service';
import { Collection } from './../../models/collection';
import { Photo } from './../../models/photo';

@Component({
  selector: 'app-photo-viewer',
  templateUrl: './photo-viewer.component.html',
  styleUrls: ['./photo-viewer.component.css'],
})
export class PhotoViewerComponent implements OnInit {
  private readonly _width = fromEvent(window, 'resize').pipe(
    map((_) => window.innerWidth)
  );

  rendition: Rendition | null;
  private _photo: Photo;

  @Input()
  set photo(photo: Photo) {
    this.rendition = null;
    this._photo = photo;
    this.pickRenditionForSize();
  }

  private _collection: Collection;
  @Input()
  set collection(collection: Collection) {
    this.rendition = null;
    this._collection = collection;
    this.pickRenditionForSize();
  }

  constructor(private readonly pathService: PathService) {}

  ngOnInit() {
    this._width.subscribe((_) => {
      this.pickRenditionForSize();
    });

    this.pickRenditionForSize();
  }

  renditionURI(): String {
    if (this.rendition === null) {
      throw 'rendition is null';
    }
    return this.pathService.rendition(this._collection, this.rendition);
  }

  pickRenditionForSize(): void {
    if (this._photo === undefined || this._collection === undefined) {
      return;
    }
    const bestFit = this._collection.renditionConfigurations.find(
      (config) => config.name === 'admin preview'
    );
    if (bestFit === undefined) {
      throw 'bestFit is undefined';
    }
    const rendition = this._photo.renditions.find(
      (rendition) => rendition.renditionConfigurationID === bestFit.id
    );
    if (rendition === undefined) {
      throw 'rendition is undefined';
    }
    this.rendition = rendition;

    console.log('Picked rendition', this.rendition);
  }
}

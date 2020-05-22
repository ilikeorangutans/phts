import { Component, OnInit } from '@angular/core';

import { switchMap, first } from 'rxjs/operators';

import { CollectionStore } from './../../stores/collection.store';
import { Collection } from './../../models/collection';
import { AlbumService } from './../../services/album.service';
import { Album } from './../../models/album';
import { Observable, BehaviorSubject } from 'rxjs';

@Component({
  selector: 'app-albums-dashboard',
  templateUrl: './albums-dashboard.component.html',
  styleUrls: ['./albums-dashboard.component.css'],
})
export class AlbumsDashboardComponent implements OnInit {
  private readonly _albums = new BehaviorSubject<Array<Album>>([]);
  albums: Observable<Array<Album>> = this._albums.asObservable();

  collection: Collection;
  album: Album = new Album();

  constructor(
    private collectionStore: CollectionStore,
    private albumService: AlbumService
  ) {}

  ngOnInit() {
    this.collectionStore.current.subscribe((collection) => {
      if (collection === null) {
        throw 'collectio is null';
      }
      this.collection = collection;
    });
    this.refreshAlbums();
  }

  refreshAlbums(): void {
    this.collectionStore.current
      .pipe(
        switchMap((collection) => {
          if (collection === null) {
            throw 'collectio is null';
          }

          return this.albumService.list(collection);
        }),
        first()
      )
      .subscribe((albums) => {
        this._albums.next(albums);
      });
  }

  onSubmit() {
    this.albumService.save(this.collection, this.album).then((_) => {
      this.refreshAlbums();
    });
    this.album = new Album();
  }
}

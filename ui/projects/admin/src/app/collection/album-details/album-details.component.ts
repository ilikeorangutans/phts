import { ActivatedRoute, Router } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { AlbumStore } from './../../stores/album.store';
import { Paginator } from '../../models/paginator';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Photo } from './../../models/photo';
import { Album } from './../../models/album';
import { Collection } from './../../models/collection';
import { AlbumService } from './../../services/album.service';
import { CollectionStore } from '../../stores/collection.store';
import { Observable, combineLatest } from 'rxjs';
import { switchMap, first } from 'rxjs/operators';

@Component({
  selector: 'app-album-details',
  templateUrl: './album-details.component.html',
  styleUrls: ['./album-details.component.css'],
})
export class AlbumDetailsComponent implements OnInit, OnDestroy {
  paginator: Paginator;

  album: Observable<Album>;
  photos: Observable<Array<Photo>>;
  collection: Collection;
  thumbnailRendition: RenditionConfiguration;

  albumStore: AlbumStore;

  constructor(
    private collectionStore: CollectionStore,
    private route: ActivatedRoute,
    private router: Router,
    private albumService: AlbumService
  ) {}

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');

    this.album = combineLatest(
      this.collectionStore.current,
      this.route.params
    ).pipe(
      switchMap(([collection, params]) => {
        const id = +params['album_id'];
        if (collection === null) {
          throw 'collection is null';
        }
        this.collection = collection;

        const renditionConfig = this.collection.renditionConfigurations.find(
          (c) => c.name === 'admin thumbnails'
        );

        if (renditionConfig === undefined) {
          throw 'rendition config is null';
        }
        this.thumbnailRendition = renditionConfig;
        return this.albumService.details(collection, id);
      })
    );

    this.album.pipe(first()).subscribe((album) => {
      this.albumStore = this.collectionStore.albumStore(album);
      this.albumStore.loadPhotos(this.paginator);
      this.photos = this.albumStore.list;
    });
  }

  ngOnDestroy(): void {}

  delete(album: Album): void {
    if (!confirm(`Delete album ${album.name}?`)) {
      return;
    }
    this.albumService
      .delete(album.collection, album)
      .subscribe((_) =>
        this.router.navigate(['collection', this.collection.slug])
      );
  }

  onPhotoClicked(photo: Photo): void {
    console.log('onPhotoClicked', photo);
    alert(`show full screen preview of photo ${photo.id}`);
  }

  setCoverPhoto(album: Album, photo: Photo): void {
    this.albumService
      .setCoverPhoto(album, photo)
      .subscribe((success) =>
        console.log('cover photo change success', success)
      );
  }

  toggleOrganizePhotos(): void {
    alert('TODO: enter organize mode');
  }

  shareAlbumDialog(): void {
    alert('TODO: share album dialog');
  }
}

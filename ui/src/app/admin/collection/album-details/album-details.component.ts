import { PhotoService } from './../../services/photo.service';
import { Photo } from './../../models/photo';
import { Observable } from 'rxjs/Observable';
import { Album } from './../../models/album';
import { Collection } from './../../models/collection';
import { Subscription } from 'rxjs/Subscription';
import { CollectionService } from './../../services/collection.service';
import { AlbumService } from './../../services/album.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Component, OnInit, OnDestroy } from '@angular/core';
import 'rxjs/observable/combineLatest';
import 'rxjs/add/operator/mergeMap';
import { combineLatest } from 'rxjs/observable/combineLatest';
import { share } from 'rxjs/operators/share';
import { Paginator } from '../../models/paginator';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';

@Component({
  selector: 'app-album-details',
  templateUrl: './album-details.component.html',
  styleUrls: ['./album-details.component.css']
})
export class AlbumDetailsComponent implements OnInit, OnDestroy {
  paginator: Paginator;

  album: Observable<Album>;
  photos: Array<Photo> = [];
  collection: Collection;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private albumService: AlbumService,
    private collectionService: CollectionService,
    private photoService: PhotoService,
    private renditionConfigService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    const x = this.collectionService.current.switchMap(collection => {
      console.log('preproccess collection');
      return Observable.fromPromise(this.renditionConfigService.forCollection(collection)).map(configs => {
        collection.renditionConfigurations = configs;
        return collection;
      });
    });
    this.album = combineLatest(x, this.route.params)
      .switchMap(([collection, params]) => {
        const id = +params['album_id'];
        this.collection = collection;
        return this.albumService.details(collection, id);
      }).pipe(share());

    this.album.switchMap(album => {
      return this.photoService.listAlbum(album.collection, album, this.paginator);
    }).subscribe(photos => this.photos = photos);
  }

  ngOnDestroy(): void {

  }

  delete(album: Album): void {
    this.albumService.delete(album.collection, album).subscribe(_ => this.router.navigate([`/admin/collection/${this.collection.slug}`]));
  }

}

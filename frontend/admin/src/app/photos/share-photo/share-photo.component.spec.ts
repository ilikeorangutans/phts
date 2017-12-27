import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SharePhotoComponent } from './share-photo.component';

describe('SharePhotoComponent', () => {
  let component: SharePhotoComponent;
  let fixture: ComponentFixture<SharePhotoComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SharePhotoComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SharePhotoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

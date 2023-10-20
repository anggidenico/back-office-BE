package controllers

import "time"

type MsPaymentChannel struct {
	PchannelKey         *string    `db:"pchannel_key"    json:"pchannel_key"      gorm:"column:pchannel_key;primary_key;autoIncrement:true"`                  /* PK table ms_payment_channel */
	PchannelCode        *string    `db:"pchannel_code"    json:"pchannel_code" gorm:"column:pchannel_code;type:varchar(50)"`                                  /* Kode PG */
	PchannelName        *string    `db:"pchannel_name"    json:"pchannel_name" gorm:"column:pchannel_name;type:varchar(150)"`                                 /* Nama yg lebih dikenal */
	SettleChannel       *string    `db:"settle_channel"    json:"settle_channel" gorm:"column:settle_channel;type:int(11)"`                                   /* Vendor/perusahaan PG penyedia payment method. lkp_group_key = 73 */
	SettlePaymentMethod *string    `db:"settle_payment_method"    json:"settle_payment_method" gorm:"column:settle_payment_method;type:int(11)"`              /* Payment method */
	MinNominalTrx       *string    `db:"min_nominal_trx"    json:"min_nominal_trx" gorm:"column:min_nominal_trx;type:decimal(18,2)"`                          /* Min. nominal transaksi yg akan dikenakan fee ini. default 0 = semua akan kena fee layanan */
	ValueType           *string    `db:"value_type"    json:"value_type" gorm:"column:value_type;type:int(11)"`                                               /* lookup value_type: FixAmount | Percentage */
	CurrencyKey         *string    `db:"currency_key"    json:"currency_key" gorm:"column:currency_key;type:int(11)"`                                         /* Mata uang Fee */
	FeeValue            *string    `db:"fee_value"    json:"fee_value" gorm:"column:fee_value;type:decimal(18,4)"`                                            /* Nilai Fee */
	HasMinMax           *string    `db:"has_min_max"    json:"has_min_max" gorm:"column:has_min_max;type:tinyint(4)"`                                         /* True: nilai min dan max harus diset */
	FeeMinValue         *string    `db:"fee_min_value"    json:"fee_min_value" gorm:"column:fee_min_value;type:decimal(18,2)"`                                /* Nilai Fee Minimum */
	FeeMaxValue         *string    `db:"fee_max_value"    json:"fee_max_value" gorm:"column:fee_max_value;type:decimal(18,2)"`                                /* Nilai Fee Maximum */
	FixedAmountFee      *string    `db:"fixed_amount_fee"    json:"fixed_amount_fee" gorm:"column:fixed_amount_fee;type:decimal(18,2)"`                       /* jika ada biaya tetap/hari, yg sifatnya selalu ada */
	FixedDmrFee         *string    `db:"fixed_dmr_fee"    json:"fixed_dmr_fee" gorm:"column:fixed_dmr_fee;type:decimal(18,2)"`                                /* jika ada fixed_dmr_fee */
	PgTnc               *string    `db:"pg_tnc"    json:"pg_tnc" gorm:"column:pg_tnc;type:text"`                                                              /* Isi TNC */
	PgRemarks           *string    `db:"pg_remarks"    json:"pg_remarks" gorm:"column:pg_remarks;type:text"`                                                  /* Remarks */
	PaymentLoginUrl     *string    `db:"payment_login_url"    json:"payment_login_url" gorm:"column:payment_login_url;type:varchar(255)"`                     /*  */
	PaymentEntryUrl     *string    `db:"payment_entry_url"    json:"payment_entry_url" gorm:"column:payment_entry_url;type:varchar(255)"`                     /*  */
	PaymentErrorUrl     *string    `db:"payment_error_url"    json:"payment_error_url" gorm:"column:payment_error_url;type:varchar(255)"`                     /*  */
	PaymentSuccessUrl   *string    `db:"payment_success_url"    json:"payment_success_url" gorm:"column:payment_success_url;type:varchar(255)"`               /*  */
	PgPrefix            *string    `db:"pg_prefix"    json:"pg_prefix" gorm:"column:pg_prefix;type:varchar(150)"`                                             /*  */
	PicName             *string    `db:"pic_name"    json:"pic_name" gorm:"column:pic_name;type:varchar(150)"`                                                /*  */
	PicPhoneNo          *string    `db:"pic_phone_no"    json:"pic_phone_no" gorm:"column:pic_phone_no;type:varchar(50)"`                                     /*  */
	PicEmailAddress     *string    `db:"pic_email_address"    json:"pic_email_address" gorm:"column:pic_email_address;type:varchar(255)"`                     /*  */
	RecOrder            uint32     `db:"rec_order"    json:"rec_order" gorm:"column:rec_order;type:int(11);default:0"`                                        /* Urutan record ditampilkan. Set value kolom ini jika ingin mengurutkan data tampil. Pada akhir setiap select query order by kolom rec_order ASC */
	RecStatus           uint8      `db:"rec_status"    json:"rec_status" gorm:"column:rec_status;type:tinyint(3) unsigned;default:1"`                         /* Status of record : 1 = active | 0 = tidak aktif | 2 = archieved | 9 = deleted. Untuk menampilan record yang aktif, pada setiap select selalu gunakan kondisi WHERE rec_status=1  */
	RecCreatedDate      *time.Time `db:"rec_created_date"    json:"rec_created_date" gorm:"column:rec_created_date;type:datetime(6);autoCreateTime:true"`     /* DateTime record diinsert. selalu isi kolom ini ketika action INSERT. tanggal diambil dari system */
	RecCreatedBy        string     `db:"rec_created_by"    json:"rec_created_by" gorm:"column:rec_created_by;type:varchar(30);default:system"`                /* Userkey/UserName yang melakukan insert. Selalu isi kolom ini ketika action INSERT. Ambil userid dari session login */
	RecModifiedDate     *time.Time `db:"rec_modified_date"    json:"rec_modified_date" gorm:"column:rec_modified_date;type:datetime(6);autoUpdateTime:milli"` /* DateTime record diupdate/ubah. Selalu isi kolom ini ketika action UPDATE. tanggal diambil dari system */
	RecModifiedBy       string     `db:"rec_modified_by"    json:"rec_modified_by" gorm:"column:rec_modified_by;type:varchar(30);default:system"`             /* User key yang melakukan perubahan pada record. Selalu isi kolom ini ketika action UPDATE. Ambil userid dari session login */
	RecImage1           string     `db:"rec_image1"    json:"rec_image1" gorm:"column:rec_image1;type:varchar(255)"`                                          /* nama icon atau image url - jika memerlukan image untuk record ini */
	RecImage2           string     `db:"rec_image2"    json:"rec_image2" gorm:"column:rec_image2;type:varchar(255)"`                                          /* nama icon ke 2 atau image url ke 2 - jika memerlukan image untuk record ini */
	RecApprovalStatus   uint8      `db:"rec_approval_status"    json:"rec_approval_status" gorm:"column:rec_approval_status;type:tinyint(3) unsigned"`        /* Approval status: 0 = Pending(Waiting) | 1 = Approved | 2 = Rejected */
	RecApprovalStage    uint32     `db:"rec_approval_stage"    json:"rec_approval_stage" gorm:"column:rec_approval_stage;type:int(11)"`                       /* Appoval stage sesuai flow approval. Nama stage lihat di workflow stage */
	RecApprovedDate     *time.Time `db:"rec_approved_date"    json:"rec_approved_date" gorm:"column:rec_approved_date;type:datetime(6)"`                      /* Tanggal approval di lakukan - ambil dari tanggal system */
	RecApprovedBy       string     `db:"rec_approved_by"    json:"rec_approved_by" gorm:"column:rec_approved_by;type:varchar(30)"`                            /* UserID yang melakukan approval - ambil dari userlogin */
	RecDeletedDate      *time.Time `db:"rec_deleted_date"    json:"rec_deleted_date" gorm:"column:rec_deleted_date;type:datetime(6)"`                         /* Tanggal record dihapus */
	RecDeletedBy        string     `db:"rec_deleted_by"    json:"rec_deleted_by" gorm:"column:rec_deleted_by;type:varchar(30)"`                               /* User yang melakukan penghapusan */
	RecAttributeId1     string     `db:"rec_attribute_id1"    json:"rec_attribute_id1" gorm:"column:rec_attribute_id1;type:varchar(50)"`                      /* Field tambahan-1 : untuk keperluan migrasi atau lainya */
	RecAttributeId2     string     `db:"rec_attribute_id2"    json:"rec_attribute_id2" gorm:"column:rec_attribute_id2;type:varchar(150)"`                     /* Field tambahan-2 : untuk keperluan migrasi atau lainya */
	RecAttributeId3     string     `db:"rec_attribute_id3"    json:"rec_attribute_id3" gorm:"column:rec_attribute_id3;type:varchar(255)"`                     /* Field tambahan-3 : untuk keperluan migrasi atau lainya */
}

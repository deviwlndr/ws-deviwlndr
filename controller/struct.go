package controller

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Karyawan struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"123456789"`
    Nama         string             `bson:"nama,omitempty" json:"nama,omitempty" example:"Tes Swagger"`
    Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty" example:"08123456789"`
    Jabatan      string             `bson:"jabatan,omitempty" json:"jabatan,omitempty" example:"Anonymous"`
    Jam_kerja    []JamKerja         `bson:"jam_kerja,omitempty" json:"jam_kerja,omitempty"`
    Hari_kerja   []string           `bson:"hari_kerja,omitempty" json:"hari_kerja,omitempty" example:"Senin,Selasa,Rabu,Kamis,Jumat,Sabtu,Minggu"`
}

type JamKerja struct {
    Durasi     int      `bson:"durasi,omitempty" json:"durasi,omitempty" example:"8"`
    Jam_masuk  string   `bson:"jam_masuk,omitempty" json:"jam_masuk,omitempty" example:"08:00"`
    Jam_keluar string   `bson:"jam_keluar,omitempty" json:"jam_keluar,omitempty" example:"16:00"`
    Gmt        int      `bson:"gmt,omitempty" json:"gmt,omitempty" example:"7"`
    Hari       []string `bson:"hari,omitempty" json:"hari,omitempty" example:"Senin,Selasa,Rabu,Kamis,Jumat,Sabtu,Minggu"`
    Shift      int      `bson:"shift,omitempty" json:"shift,omitempty" example:"2"`
    Piket_tim  string   `bson:"piket_tim,omitempty" json:"piket_tim,omitempty" example:"Piket Z"`
}

type Presensi struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"123456789"`
    Longitude    float64            `bson:"longitude,omitempty" json:"longitude,omitempty" example:"123.11"`
    Latitude     float64            `bson:"latitude,omitempty" json:"latitude,omitempty" example:"123.12"`
    Location     string             `bson:"location,omitempty" json:"location,omitempty" example:"Bandung"`
    Phone_number string             `bson:"phone_number,omitempty" json:"phone_number,omitempty" example:"08123456789"`
    Datetime     primitive.DateTime `bson:"datetime,omitempty" json:"datetime,omitempty" swaggertype:"string" example:"2024-09-01T00:00:00Z" format:"date-time"`
    Checkin      string             `bson:"checkin,omitempty" json:"checkin,omitempty" example:"Masuk"`
    Biodata      Karyawan           `bson:"biodata,omitempty" json:"biodata,omitempty"`
}
